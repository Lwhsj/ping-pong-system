package service

import (
	"pingpong-backend/internal/model"

	"gorm.io/gorm"
)

type PlayerService struct {
	db *gorm.DB
}

func NewPlayerService(db *gorm.DB) *PlayerService {
	return &PlayerService{db: db}
}

func (s *PlayerService) GetAllPlayers() ([]model.Player, error) {
	var players []model.Player
	err := s.db.Find(&players).Error
	return players, err
}
