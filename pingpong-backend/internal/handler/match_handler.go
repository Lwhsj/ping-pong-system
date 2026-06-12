package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pingpong-backend/internal/dto"
	"pingpong-backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MatchHandler struct {
	matchService *service.MatchService
}

func NewMatchHandler(matchService *service.MatchService) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

func (h *MatchHandler) StartMatch(c *gin.Context) {
	var request dto.MatchStartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	match, err := h.matchService.StartMatch(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) GetMatch(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	match, err := h.matchService.GetMatch(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) FinishMatch(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	match, err := h.matchService.FinishMatch(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) CurrentScore(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	response, err := h.matchService.CurrentScore(id)
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

func (h *MatchHandler) MatchHistory(c *gin.Context) {
	dateFrom, ok := parseOptionalDate(c, "dateFrom")
	if !ok {
		return
	}
	dateTo, ok := parseOptionalDate(c, "dateTo")
	if !ok {
		return
	}

	matches, err := h.matchService.MatchHistory(dateFrom, dateTo, c.Query("player"), c.Query("status"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

func (h *MatchHandler) MatchDetail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	rallies, err := h.matchService.MatchDetail(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rallies)
}

func (h *MatchHandler) MatchStats(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	stats, err := h.matchService.MatchStats(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *MatchHandler) ExportMatch(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	buffer := bytes.NewBuffer(nil)
	if err := h.matchService.ExportMatch(id, buffer); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("match_%d.xlsx", id)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename=`+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}

func parseUintParam(c *gin.Context, name string) (uint64, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return value, true
}

func parseOptionalDate(c *gin.Context, name string) (*time.Time, bool) {
	value := c.Query(name)
	if value == "" {
		return nil, true
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad date format"})
		return nil, false
	}
	return &parsed, true
}
