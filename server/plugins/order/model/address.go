package model

import "github.com/ijry/lyshop/model"

type Address struct {
	model.Base
	UserID    uint64 `gorm:"not null;index"    json:"user_id"`
	Name      string `gorm:"size:64;not null"  json:"name"`
	Phone     string `gorm:"size:20;not null"  json:"phone"`
	Province  string `gorm:"size:32"           json:"province"`
	City      string `gorm:"size:32"           json:"city"`
	District  string `gorm:"size:32"           json:"district"`
	Detail    string `gorm:"size:255"          json:"detail"`
	IsDefault int8   `gorm:"not null;default:0" json:"is_default"`
}
