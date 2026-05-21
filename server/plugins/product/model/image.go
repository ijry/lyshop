package model

import "github.com/ijry/lyshop/model"

type ProductImage struct {
	model.Base
	ProductID uint64 `gorm:"not null;index" json:"product_id"`
	URL       string `gorm:"size:500;not null" json:"url"`
	Sort      int    `gorm:"not null;default:0" json:"sort"`
}
