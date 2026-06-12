package dto

type AgentAnalyzeMatchRequest struct {
	Question string `json:"question"`
}

type AgentChatRequest struct {
	MatchID  uint64 `json:"match_id" binding:"required"`
	Question string `json:"question" binding:"required"`
}

type AgentMatchAnalysisResponse struct {
	MatchID             uint64            `json:"match_id"`
	Summary             string            `json:"summary"`
	Strengths           []string          `json:"strengths"`
	Weaknesses          []string          `json:"weaknesses"`
	KeyMoments          []AgentKeyMoment  `json:"key_moments"`
	TrainingSuggestions []string          `json:"training_suggestions"`
	RawText             string            `json:"raw_text,omitempty"`
	Metadata            map[string]string `json:"metadata,omitempty"`
}

type AgentKeyMoment struct {
	RallyNumber int    `json:"rally_number"`
	Reason      string `json:"reason"`
}

type AgentChatResponse struct {
	MatchID  uint64 `json:"match_id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
