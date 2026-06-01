package model

import (
	"encoding/json"

	"github.com/ijry/lyshop/model"
)

type SpecTemplateAttr struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type SpecTemplate struct {
	model.Base
	Name        string          `gorm:"size:128;not null;index" json:"name"`
	CategoryIDs json.RawMessage `gorm:"type:json"               json:"category_ids"` // []uint64
	Attrs       json.RawMessage `gorm:"type:json"               json:"attrs"`        // []SpecTemplateAttr
	Status      int8            `gorm:"not null;default:1;index" json:"status"`
	Sort        int             `gorm:"not null;default:0;index" json:"sort"`
}
