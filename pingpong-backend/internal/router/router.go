package router

import (
	"log"

	"pingpong-backend/internal/config"
	"pingpong-backend/internal/handler"
	"pingpong-backend/internal/middleware"
	"pingpong-backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(db *gorm.DB, cfg config.Config) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	playerService := service.NewPlayerService(db)
	matchService := service.NewMatchService(db)
	rallyService := service.NewRallyService(db)
	llmClient, agentErr := service.NewLLMClient(cfg)
	agentService := service.NewAgentService(db, llmClient)

	playerHandler := handler.NewPlayerHandler(playerService)
	matchHandler := handler.NewMatchHandler(matchService)
	rallyHandler := handler.NewRallyHandler(rallyService)
	agentHandler := handler.NewAgentHandler(agentService, agentErr)
	videoHandler, err := handler.NewVideoHandler(cfg)
	if err != nil {
		log.Fatalf("create video handler: %v", err)
	}

	api := r.Group("/api")
	{
		api.GET("/players", playerHandler.GetAllPlayers)

		api.POST("/match/start", matchHandler.StartMatch)
		api.GET("/match/:id", matchHandler.GetMatch)
		api.POST("/match/:id/finish", matchHandler.FinishMatch)
		api.GET("/match/:id/current", matchHandler.CurrentScore)
		api.GET("/matches", matchHandler.MatchHistory)
		api.GET("/match/:id/detail", matchHandler.MatchDetail)
		api.GET("/match/:id/stats", matchHandler.MatchStats)
		api.GET("/match/:id/export", matchHandler.ExportMatch)

		api.POST("/rally", rallyHandler.SaveRally)

		api.POST("/agent/match/:id/analyze", agentHandler.AnalyzeMatch)
		api.POST("/agent/chat", agentHandler.Chat)

		api.POST("/upload/video", videoHandler.UploadVideo)
		api.GET("/video/:fileName", videoHandler.StreamVideo)
	}

	return r
}
