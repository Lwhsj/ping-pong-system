package dto

type MatchStartRequest struct {
	Player1ID   uint64  `json:"player1_id" binding:"required"`
	Player2ID   uint64  `json:"player2_id" binding:"required"`
	FirstServer *uint64 `json:"first_server"`
}

type MatchCurrentScoreResponse struct {
	MatchID     uint64 `json:"match_id"`
	RallyNumber int    `json:"rally_number"`
	ScoreP1     int    `json:"score_p1"`
	ScoreP2     int    `json:"score_p2"`
	Server      string `json:"server"`
	Player1Name string `json:"player1_name,omitempty"`
	Player2Name string `json:"player2_name,omitempty"`
}

type MatchStatsResponse struct {
	ServeSuccessRate map[string]float64 `json:"serve_success_rate"`
	ConsecutiveScore map[string]int     `json:"consecutive_score"`
	AverageRallyTime float64            `json:"average_rally_time"`
}
