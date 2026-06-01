package handler

import (
	"net/http"

	"pingpong-backend/internal/dto"
	"pingpong-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type RallyHandler struct {
	rallyService *service.RallyService
}

func NewRallyHandler(rallyService *service.RallyService) *RallyHandler {
	return &RallyHandler{rallyService: rallyService}
}

func (h *RallyHandler) SaveRally(c *gin.Context) {
	var request dto.RallyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rally, err := h.rallyService.SaveRally(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rally)
}
