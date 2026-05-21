package model

import (
	"github.com/ijry/lyshop/model"
)

type Product struct {
	model.Base
	MerchantID  uint64  `gorm:"not null;default:0;index"  json:"merchant_id"`
	CategoryID  uint64  `gorm:"not null;index"             json:"category_id"`
	Title       string  `gorm:"size:255;not null"          json:"title"`
	Subtitle    string  `gorm:"size:255"                   json:"subtitle"`
	Cover       string  `gorm:"size:500"                   json:"cover"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginPrice float64 `gorm:"type:decimal(10,2)"         json:"origin_price"`
	Stock       int     `gorm:"not null;default:0"         json:"stock"`
	Sales       int     `gorm:"not null;default:0"         json:"sales"`
	Status      int8    `gorm:"not null;default:1"         json:"status"` // 1=上架 0=下架
	Sort        int     `gorm:"not null;default:0"         json:"sort"`
	Detail      string  `gorm:"type:longtext"              json:"detail"`
}
