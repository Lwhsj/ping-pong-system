package service

import (
	"pingpong-backend/internal/dto"
	"pingpong-backend/internal/model"

	"gorm.io/gorm"
)

type RallyService struct {
	db *gorm.DB
}

func NewRallyService(db *gorm.DB) *RallyService {
	return &RallyService{db: db}
}

func (s *RallyService) SaveRally(request dto.RallyRequest) (model.Rally, error) {
	var rallies []model.Rally
	if err := s.db.Where("match_id = ?", request.MatchID).Find(&rallies).Error; err != nil {
		return model.Rally{}, err
	}

	scoreP1, scoreP2 := CalculateNextScore(rallies, request.Scorer)
	rally := model.Rally{
		MatchID:     request.MatchID,
		RallyNumber: request.RallyNumber,
		Scorer:      request.Scorer,
		Server:      request.Server,
		Timestamp:   request.Timestamp.Time,
		VideoFile:   request.VideoFile,
		ScoreP1:     &scoreP1,
		ScoreP2:     &scoreP2,
	}

	if err := s.db.Create(&rally).Error; err != nil {
		return model.Rally{}, err
	}
	return rally, nil
}

func (s *RallyService) GetRalliesByMatchID(matchID uint64) ([]model.Rally, error) {
	var rallies []model.Rally
	err := s.db.Where("match_id = ?", matchID).Order("rally_number ASC").Find(&rallies).Error
	return rallies, err
}
