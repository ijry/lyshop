package model

import (
	"encoding/json"

	"github.com/ijry/lyshop/model"
)

// SkuAttr represents one attribute key-value pair, e.g. {"name":"颜色","value":"红色"}
type SkuAttr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ProductSku struct {
	model.Base
	ProductID uint64          `gorm:"not null;index"             json:"product_id"`
	Attrs     json.RawMessage `gorm:"type:json"                  json:"attrs"` // []SkuAttr
	Price     float64         `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock     int             `gorm:"not null;default:0"         json:"stock"`
	SkuCode   string          `gorm:"size:128"                   json:"sku_code"`
}
