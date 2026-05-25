package model

import "github.com/ijry/lyshop/model"

type ProductFavorite struct {
	model.Base
	UserID    uint64 `gorm:"not null;index;uniqueIndex:uniq_user_product" json:"user_id"`
	ProductID uint64 `gorm:"not null;index;uniqueIndex:uniq_user_product" json:"product_id"`
}
