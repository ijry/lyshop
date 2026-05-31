// Package model provides base types shared across domain models.
package model

import "time"

// Base provides common fields for all models.
type Base struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Time is an alias for *time.Time for use in optional timestamp fields.
type Time = time.Time
