package model

import "time"

type Rally struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	MatchID     uint64    `gorm:"column:match_id;not null" json:"match_id"`
	RallyNumber int       `gorm:"column:rally_number;not null" json:"rally_number"`
	Scorer      string    `gorm:"column:scorer;size:20;not null" json:"scorer"`
	Server      string    `gorm:"column:server;size:20;not null" json:"server"`
	Timestamp   time.Time `gorm:"column:timestamp;not null" json:"timestamp"`
	VideoFile   string    `gorm:"column:video_file;size:100;not null" json:"video_file"`
	ScoreP1     *int      `gorm:"column:score_p1" json:"score_p1"`
	ScoreP2     *int      `gorm:"column:score_p2" json:"score_p2"`
}

func (Rally) TableName() string {
	return "rally"
}
