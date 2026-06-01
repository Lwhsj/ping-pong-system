package handler

import (
	"net/http"

	"pingpong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService *service.PlayerService
}

func NewPlayerHandler(playerService *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{playerService: playerService}
}

func (h *PlayerHandler) GetAllPlayers(c *gin.Context) {
	players, err := h.playerService.GetAllPlayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}
