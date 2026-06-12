package service

import (
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"pingpong-backend/internal/dto"
	"pingpong-backend/internal/model"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type MatchService struct {
	db *gorm.DB
}

func NewMatchService(db *gorm.DB) *MatchService {
	return &MatchService{db: db}
}

func (s *MatchService) StartMatch(request dto.MatchStartRequest) (model.Match, error) {
	match := model.Match{
		Date:        model.TodayDate(),
		Player1ID:   request.Player1ID,
		Player2ID:   request.Player2ID,
		Status:      "started",
		StartedAt:   time.Now(),
		FirstServer: request.FirstServer,
	}
	if err := s.db.Create(&match).Error; err != nil {
		return model.Match{}, err
	}
	return match, nil
}

func (s *MatchService) GetMatch(id uint64) (model.Match, error) {
	var match model.Match
	if err := s.db.First(&match, id).Error; err != nil {
		return model.Match{}, err
	}
	match.Player1Name, match.Player2Name = s.playerNames(match)
	return match, nil
}

func (s *MatchService) FinishMatch(id uint64) (model.Match, error) {
	var match model.Match
	if err := s.db.First(&match, id).Error; err != nil {
		return model.Match{}, err
	}
	match.Status = "finished"
	if err := s.db.Save(&match).Error; err != nil {
		return model.Match{}, err
	}
	return match, nil
}

func (s *MatchService) CurrentScore(id uint64) (dto.MatchCurrentScoreResponse, error) {
	var rallies []model.Rally
	if err := s.db.Where("match_id = ?", id).Find(&rallies).Error; err != nil {
		return dto.MatchCurrentScoreResponse{}, err
	}

	scoreP1, scoreP2 := LatestScore(rallies)
	rallyNumber := 0
	if len(rallies) > 0 {
		rallyNumber = LatestRally(rallies).RallyNumber
	}

	var match model.Match
	err := s.db.First(&match, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.MatchCurrentScoreResponse{}, err
	}
	if err != nil {
		return dto.MatchCurrentScoreResponse{}, err
	}

	p1ID := match.Player1ID
	p2ID := match.Player2ID
	response := dto.MatchCurrentScoreResponse{
		MatchID:     id,
		RallyNumber: rallyNumber,
		ScoreP1:     scoreP1,
		ScoreP2:     scoreP2,
		Server:      CalculateServer(scoreP1, scoreP2, match.FirstServer, &p1ID, &p2ID),
	}
	response.Player1Name, response.Player2Name = s.playerNames(match)
	return response, nil
}

func (s *MatchService) MatchHistory(dateFrom, dateTo *time.Time, playerName, status string) ([]model.Match, error) {
	var matches []model.Match
	query := s.db.Model(&model.Match{})
	if dateFrom != nil {
		query = query.Where("date >= ?", dateFrom.Format("2006-01-02"))
	}
	if dateTo != nil {
		query = query.Where("date <= ?", dateTo.Format("2006-01-02"))
	}
	if status != "" {
		query = query.Where("LOWER(status) = LOWER(?)", status)
	}
	if playerName != "" {
		query = query.Joins("LEFT JOIN player p1 ON p1.id = matches.player1_id").
			Joins("LEFT JOIN player p2 ON p2.id = matches.player2_id").
			Where("p1.name LIKE ? OR p2.name LIKE ?", "%"+playerName+"%", "%"+playerName+"%")
	}
	if err := query.Find(&matches).Error; err != nil {
		return nil, err
	}

	for i := range matches {
		matches[i].Player1Name, matches[i].Player2Name = s.playerNames(matches[i])
	}
	return matches, nil
}

func (s *MatchService) MatchDetail(id uint64) ([]model.Rally, error) {
	var rallies []model.Rally
	err := s.db.Where("match_id = ?", id).Order("rally_number ASC").Find(&rallies).Error
	return rallies, err
}

func (s *MatchService) MatchStats(id uint64) (dto.MatchStatsResponse, error) {
	var rallies []model.Rally
	if err := s.db.Where("match_id = ?", id).Find(&rallies).Error; err != nil {
		return dto.MatchStatsResponse{}, err
	}

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

	return dto.MatchStatsResponse{
		ServeSuccessRate: map[string]float64{
			"player1": rate(p1ServeWin, p1ServeCount),
			"player2": rate(p2ServeWin, p2ServeCount),
		},
		ConsecutiveScore: map[string]int{
			"player1": MaxConsecutiveScore(append([]model.Rally(nil), rallies...), "player1"),
			"player2": MaxConsecutiveScore(append([]model.Rally(nil), rallies...), "player2"),
		},
		AverageRallyTime: 15.2,
	}, nil
}

func (s *MatchService) ExportMatch(id uint64, out io.Writer) error {
	var match model.Match
	if err := s.db.First(&match, id).Error; err != nil {
		return err
	}
	var rallies []model.Rally
	if err := s.db.Where("match_id = ?", id).Order("rally_number ASC").Find(&rallies).Error; err != nil {
		return err
	}

	player1Name, player2Name := s.playerNames(match)
	if player1Name == "" {
		player1Name = "Unknown"
	}
	if player2Name == "" {
		player2Name = "Unknown"
	}

	file := excelize.NewFile()
	defer file.Close()

	defaultSheet := file.GetSheetName(0)
	sheet := "姣旇禌鏁版嵁"
	if err := file.SetSheetName(defaultSheet, sheet); err != nil {
		return err
	}

	columns := []string{"鍥炲悎鏁?", "寰楀垎鑰?", "姣斿垎 (P1-P2)", "鍙戠悆鏂?", "鏃堕棿"}
	for idx, column := range columns {
		cell, err := excelize.CoordinatesToCellName(idx+1, 1)
		if err != nil {
			return err
		}
		if err := file.SetCellValue(sheet, cell, column); err != nil {
			return err
		}
	}

	sort.Slice(rallies, func(i, j int) bool {
		return rallies[i].RallyNumber < rallies[j].RallyNumber
	})
	for idx, rally := range rallies {
		row := idx + 2
		values := []any{
			rally.RallyNumber,
			playerLabel(rally.Scorer, player1Name, player2Name),
			scoreLabel(rally.ScoreP1, rally.ScoreP2),
			playerLabel(rally.Server, player1Name, player2Name),
			rally.Timestamp.Format("2006-01-02 15:04:05"),
		}
		for col, value := range values {
			cell, err := excelize.CoordinatesToCellName(col+1, row)
			if err != nil {
				return err
			}
			if err := file.SetCellValue(sheet, cell, value); err != nil {
				return err
			}
		}
	}
	return file.Write(out)
}

func (s *MatchService) playerNames(match model.Match) (string, string) {
	var p1, p2 model.Player
	player1Name, player2Name := "", ""
	if err := s.db.First(&p1, match.Player1ID).Error; err == nil {
		player1Name = p1.Name
	}
	if err := s.db.First(&p2, match.Player2ID).Error; err == nil {
		player2Name = p2.Name
	}
	return player1Name, player2Name
}

func rate(wins, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(wins) / float64(total)
}

func playerLabel(value, player1Name, player2Name string) string {
	if strings.EqualFold(value, "player1") {
		return player1Name
	}
	if strings.EqualFold(value, "player2") {
		return player2Name
	}
	return value
}

func scoreLabel(scoreP1, scoreP2 *int) string {
	p1, p2 := 0, 0
	if scoreP1 != nil {
		p1 = *scoreP1
	}
	if scoreP2 != nil {
		p2 = *scoreP2
	}
	return strings.Join([]string{itoa(p1), itoa(p2)}, "-")
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
