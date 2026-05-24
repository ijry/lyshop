package model

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/model"
)

// DecorPage stores the component configuration for a decorated page.
type DecorPage struct {
	model.Base
	MerchantID  uint64          `gorm:"not null;default:0;uniqueIndex:uk_merchant_page_variant" json:"merchant_id"`
	PageKey     string          `gorm:"size:64;not null;uniqueIndex:uk_merchant_page_variant"   json:"page_key"` // index|category|...
	VariantKey  string          `gorm:"size:64;not null;default:'default';uniqueIndex:uk_merchant_page_variant" json:"variant_key"`
	VariantName string          `gorm:"size:128;not null;default:'默认副本'"                     json:"variant_name"`
	Components  json.RawMessage `gorm:"type:json;not null"                                        json:"components"`
	IsPublished bool            `gorm:"not null;default:false;index:idx_decor_publish"            json:"is_published"`
	PublishedAt *time.Time      `json:"published_at"`
}
