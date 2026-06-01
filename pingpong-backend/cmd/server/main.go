package main

import (
	"log"

	"pingpong-backend/internal/config"
	"pingpong-backend/internal/database"
	"pingpong-backend/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	r := router.New(db, cfg)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
