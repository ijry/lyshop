package model

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/model"
)

// DecorPage stores the component configuration for a decorated page.
type DecorPage struct {
	model.Base
	MerchantID  uint64          `gorm:"not null;default:0;uniqueIndex:uk_merchant_page" json:"merchant_id"`
	PageKey     string          `gorm:"size:64;not null;uniqueIndex:uk_merchant_page"   json:"page_key"` // index|category|...
	Components  json.RawMessage `gorm:"type:json;not null"                              json:"components"`
	PublishedAt *time.Time      `json:"published_at"`
}
