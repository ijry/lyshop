package model

import "encoding/json"

// Role defines a set of permission strings for admins.
type Role struct {
	Base
	Name        string          `gorm:"size:64;not null" json:"name"`
	Permissions json.RawMessage `gorm:"type:json"        json:"permissions"` // []string
}
