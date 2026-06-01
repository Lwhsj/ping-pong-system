package dto

import (
	"strings"
	"time"
)

type LocalDateTime struct {
	time.Time
}

func (ldt *LocalDateTime) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`)
	if value == "" || value == "null" {
		ldt.Time = time.Time{}
		return nil
	}

	layouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
	}
	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			ldt.Time = parsed
			return nil
		}
		lastErr = err
	}
	return lastErr
}

type RallyRequest struct {
	MatchID     uint64        `json:"match_id" binding:"required"`
	SetNumber   *int          `json:"set_number"`
	RallyNumber int           `json:"rally_number" binding:"required"`
	Scorer      string        `json:"scorer" binding:"required"`
	Server      string        `json:"server" binding:"required"`
	Timestamp   LocalDateTime `json:"timestamp" binding:"required"`
	VideoFile   string        `json:"video_file" binding:"required"`
}
