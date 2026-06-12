package handler

import (
	"errors"
	"io"
	"net/http"

	"pingpong-backend/internal/dto"
	"pingpong-backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AgentHandler struct {
	agentService *service.AgentService
	initErr      error
}

func NewAgentHandler(agentService *service.AgentService, initErr error) *AgentHandler {
	return &AgentHandler{agentService: agentService, initErr: initErr}
}

func (h *AgentHandler) AnalyzeMatch(c *gin.Context) {
	if h.unavailable(c) {
		return
	}

	matchID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var request dto.AgentAnalyzeMatchRequest
	if err := c.ShouldBindJSON(&request); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.agentService.AnalyzeMatch(c.Request.Context(), matchID, request.Question)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *AgentHandler) Chat(c *gin.Context) {
	if h.unavailable(c) {
		return
	}

	var request dto.AgentChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.agentService.Chat(c.Request.Context(), request)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *AgentHandler) unavailable(c *gin.Context) bool {
	if h.initErr == nil {
		return false
	}
	if errors.Is(h.initErr, service.ErrAgentDisabled) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Agent is disabled"})
		return true
	}
	if errors.Is(h.initErr, service.ErrLLMNotConfigured) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "LLM API key is not configured"})
		return true
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": h.initErr.Error()})
	return true
}
