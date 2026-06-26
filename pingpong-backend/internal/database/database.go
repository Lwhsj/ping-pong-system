package database

import (
	"pingpong-backend/internal/config"
	"pingpong-backend/internal/model"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	switch strings.ToLower(cfg.DBDriver) {
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		if err := autoMigrateAndSeed(db); err != nil {
			return nil, err
		}
		return db, nil
	default:
		return gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	}
}

func autoMigrateAndSeed(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Player{}, &model.Match{}, &model.Rally{}); err != nil {
		return err
	}

	var playerCount int64
	if err := db.Model(&model.Player{}).Count(&playerCount).Error; err != nil {
		return err
	}
	if playerCount > 0 {
		return nil
	}

	seedPlayers := []model.Player{
		{Name: "Ma Long", Sex: "M", Age: 35},
		{Name: "Fan Zhendong", Sex: "M", Age: 27},
		{Name: "Sun Yingsha", Sex: "F", Age: 24},
		{Name: "Chen Meng", Sex: "F", Age: 30},
	}
	return db.Create(&seedPlayers).Error
}
