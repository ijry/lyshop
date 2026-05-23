package model

import "github.com/ijry/lyshop/model"

// CheckinRule defines how many points a user earns per consecutive day.
// day=1 means first day, day=7 means 7th consecutive day, day=0 is the default fallback.
type CheckinRule struct {
	model.Base
	Day    int `gorm:"not null;uniqueIndex" json:"day"`    // 0=default, 1-7=consecutive day
	Points int `gorm:"not null"             json:"points"` // points awarded
}

// CheckinLog records each daily check-in.
type CheckinLog struct {
	model.Base
	UserID         uint64 `gorm:"not null;index"                          json:"user_id"`
	Date           string `gorm:"size:10;not null;uniqueIndex:uk_user_date" json:"date"` // YYYY-MM-DD
	ConsecutiveDays int   `gorm:"not null"                                json:"consecutive_days"`
	Points         int    `gorm:"not null"                                json:"points"`
}

func (CheckinLog) TableName() string { return "checkin_logs" }
