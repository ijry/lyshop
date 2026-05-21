package model

import "github.com/ijry/lyshop/model"

type ProductCategory struct {
	model.Base
	ParentID uint64 `gorm:"not null;default:0;index"  json:"parent_id"`
	Name     string `gorm:"size:64;not null"           json:"name"`
	Icon     string `gorm:"size:500"                   json:"icon"`
	Sort     int    `gorm:"not null;default:0"         json:"sort"`
	Status   int8   `gorm:"not null;default:1"         json:"status"`
}
