package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func TodayDate() Date {
	now := time.Now()
	return Date{Time: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())}
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

func (d *Date) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		d.Time = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
		return nil
	case []byte:
		return d.scanString(string(v))
	case string:
		return d.scanString(v)
	case nil:
		d.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("scan Date from %T", value)
	}
}

func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Format("2006-01-02"), nil
}

func (d *Date) scanString(value string) error {
	parsed, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

type Match struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Date        Date      `gorm:"column:date;type:date;not null" json:"date"`
	Player1ID   uint64    `gorm:"column:player1_id;not null" json:"player1_id"`
	Player2ID   uint64    `gorm:"column:player2_id;not null" json:"player2_id"`
	Status      string    `gorm:"column:status;size:10;not null" json:"status"`
	StartedAt   time.Time `gorm:"column:started_at;not null" json:"started_at"`
	FirstServer *uint64   `gorm:"column:first_server" json:"first_server"`
	Player1Name string    `gorm:"-" json:"player1_name,omitempty"`
	Player2Name string    `gorm:"-" json:"player2_name,omitempty"`
}

func (Match) TableName() string {
	return "matches"
}
