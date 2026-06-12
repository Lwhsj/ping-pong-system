package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"pingpong-backend/internal/dto"
	"pingpong-backend/internal/model"

	"gorm.io/gorm"
)

type AgentService struct {
	db        *gorm.DB
	llmClient *LLMClient
}

type agentMatchContext struct {
	Match   model.Match             `json:"match"`
	Players map[string]model.Player `json:"players"`
	Score   agentScoreSummary       `json:"score"`
	Stats   agentStatsSummary       `json:"stats"`
	Rallies []agentRallyContext     `json:"rallies"`
}

type agentStatsSummary struct {
	ServeSuccessRate map[string]float64 `json:"serve_success_rate"`
	ConsecutiveScore map[string]int     `json:"consecutive_score"`
}

type agentScoreSummary struct {
	Player1 int    `json:"player1"`
	Player2 int    `json:"player2"`
	Leader  string `json:"leader,omitempty"`
	Winner  string `json:"winner,omitempty"`
}

type agentRallyContext struct {
	RallyNumber int    `json:"rally_number"`
	Scorer      string `json:"scorer"`
	Server      string `json:"server"`
	Score       string `json:"score"`
	VideoFile   string `json:"video_file,omitempty"`
}

func NewAgentService(db *gorm.DB, llmClient *LLMClient) *AgentService {
	return &AgentService{db: db, llmClient: llmClient}
}

func (s *AgentService) AnalyzeMatch(ctx context.Context, matchID uint64, question string) (dto.AgentMatchAnalysisResponse, error) {
	matchContext, err := s.buildMatchContext(matchID)
	if err != nil {
		return dto.AgentMatchAnalysisResponse{}, err
	}

	contextJSON, err := json.Marshal(matchContext)
	if err != nil {
		return dto.AgentMatchAnalysisResponse{}, err
	}

	userPrompt := strings.TrimSpace(fmt.Sprintf(`比赛数据如下：
%s

用户关注点：%s

请基于比赛数据生成复盘。`, string(contextJSON), defaultQuestion(question, "请从比分走势、发球表现、连续得失分和训练建议四个角度复盘。")))

	content, err := s.llmClient.Chat(ctx, matchAnalysisSystemPrompt(), userPrompt, true)
	if err != nil {
		return dto.AgentMatchAnalysisResponse{}, err
	}

	response := dto.AgentMatchAnalysisResponse{MatchID: matchID}
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		response.RawText = content
		response.Metadata = map[string]string{"parse_error": err.Error()}
		return response, nil
	}
	response.MatchID = matchID
	return response, nil
}

func (s *AgentService) Chat(ctx context.Context, request dto.AgentChatRequest) (dto.AgentChatResponse, error) {
	matchContext, err := s.buildMatchContext(request.MatchID)
	if err != nil {
		return dto.AgentChatResponse{}, err
	}

	contextJSON, err := json.Marshal(matchContext)
	if err != nil {
		return dto.AgentChatResponse{}, err
	}

	userPrompt := strings.TrimSpace(fmt.Sprintf(`比赛数据如下：
%s

用户问题：%s`, string(contextJSON), request.Question))

	answer, err := s.llmClient.Chat(ctx, matchChatSystemPrompt(), userPrompt, false)
	if err != nil {
		return dto.AgentChatResponse{}, err
	}

	return dto.AgentChatResponse{
		MatchID:  request.MatchID,
		Question: request.Question,
		Answer:   strings.TrimSpace(answer),
	}, nil
}

func (s *AgentService) buildMatchContext(matchID uint64) (agentMatchContext, error) {
	var match model.Match
	if err := s.db.First(&match, matchID).Error; err != nil {
		return agentMatchContext{}, err
	}

	players, err := s.loadPlayers(match)
	if err != nil {
		return agentMatchContext{}, err
	}

	var rallies []model.Rally
	if err := s.db.Where("match_id = ?", matchID).Order("rally_number ASC").Find(&rallies).Error; err != nil {
		return agentMatchContext{}, err
	}

	stats := calculateAgentStats(rallies)
	score := summarizeScore(rallies, match.Status)
	match.Player1Name = players["player1"].Name
	match.Player2Name = players["player2"].Name

	return agentMatchContext{
		Match:   match,
		Players: players,
		Score:   score,
		Stats:   stats,
		Rallies: buildRallyContext(rallies),
	}, nil
}

func (s *AgentService) loadPlayers(match model.Match) (map[string]model.Player, error) {
	players := map[string]model.Player{
		"player1": {ID: match.Player1ID},
		"player2": {ID: match.Player2ID},
	}

	var player1 model.Player
	if err := s.db.First(&player1, match.Player1ID).Error; err != nil {
		return nil, err
	}
	var player2 model.Player
	if err := s.db.First(&player2, match.Player2ID).Error; err != nil {
		return nil, err
	}
	players["player1"] = player1
	players["player2"] = player2
	return players, nil
}

func calculateAgentStats(rallies []model.Rally) agentStatsSummary {
	p1ServeCount, p1ServeWin := 0, 0
	p2ServeCount, p2ServeWin := 0, 0
	for _, rally := range rallies {
		if strings.EqualFold(rally.Server, "player1") {
			p1ServeCount++
			if strings.EqualFold(rally.Scorer, "player1") {
				p1ServeWin++
			}
		}
		if strings.EqualFold(rally.Server, "player2") {
			p2ServeCount++
			if strings.EqualFold(rally.Scorer, "player2") {
				p2ServeWin++
			}
		}
	}

	return agentStatsSummary{
		ServeSuccessRate: map[string]float64{
			"player1": rate(p1ServeWin, p1ServeCount),
			"player2": rate(p2ServeWin, p2ServeCount),
		},
		ConsecutiveScore: map[string]int{
			"player1": MaxConsecutiveScore(append([]model.Rally(nil), rallies...), "player1"),
			"player2": MaxConsecutiveScore(append([]model.Rally(nil), rallies...), "player2"),
		},
	}
}

func summarizeScore(rallies []model.Rally, matchStatus string) agentScoreSummary {
	if len(rallies) == 0 {
		return agentScoreSummary{}
	}

	sortedRallies := append([]model.Rally(nil), rallies...)
	sort.Slice(sortedRallies, func(i, j int) bool {
		return sortedRallies[i].RallyNumber < sortedRallies[j].RallyNumber
	})
	latest := sortedRallies[len(sortedRallies)-1]

	scoreP1, scoreP2 := 0, 0
	if latest.ScoreP1 != nil {
		scoreP1 = *latest.ScoreP1
	}
	if latest.ScoreP2 != nil {
		scoreP2 = *latest.ScoreP2
	}

	leader := ""
	if scoreP1 > scoreP2 {
		leader = "player1"
	}
	if scoreP2 > scoreP1 {
		leader = "player2"
	}

	winner := ""
	if strings.EqualFold(matchStatus, "finished") {
		winner = leader
	}
	return agentScoreSummary{Player1: scoreP1, Player2: scoreP2, Leader: leader, Winner: winner}
}

func buildRallyContext(rallies []model.Rally) []agentRallyContext {
	result := make([]agentRallyContext, 0, len(rallies))
	for _, rally := range rallies {
		scoreP1, scoreP2 := 0, 0
		if rally.ScoreP1 != nil {
			scoreP1 = *rally.ScoreP1
		}
		if rally.ScoreP2 != nil {
			scoreP2 = *rally.ScoreP2
		}
		result = append(result, agentRallyContext{
			RallyNumber: rally.RallyNumber,
			Scorer:      rally.Scorer,
			Server:      rally.Server,
			Score:       fmt.Sprintf("%d-%d", scoreP1, scoreP2),
			VideoFile:   rally.VideoFile,
		})
	}
	return result
}

func defaultQuestion(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func matchAnalysisSystemPrompt() string {
	return `你是专业乒乓球比赛复盘 Agent。你只能依据传入的结构化比赛数据分析，不要编造视频中没有提供的技术细节。
请使用中文输出一个 JSON object，字段必须为：
summary: string，简短总结；
strengths: string[]，优势；
weaknesses: string[]，问题；
key_moments: array，每项包含 rally_number:number 和 reason:string；
training_suggestions: string[]，可执行训练建议。
不要输出 Markdown，不要输出 JSON 以外的内容。`
}

func matchChatSystemPrompt() string {
	return `你是专业乒乓球比赛问答 Agent。你只能依据传入的结构化比赛数据回答。
如果数据不足以判断，请明确说明缺少哪些数据。回答要中文、简洁、可执行。`
}
