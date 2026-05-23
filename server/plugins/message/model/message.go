package model

import "github.com/ijry/lyshop/model"

const (
	GroupSystem    = "system"
	GroupOrder     = "order"
	GroupMarketing = "marketing"
	GroupIM        = "im"
)

type Message struct {
	model.Base
	UserID  uint64 `gorm:"not null;index"    json:"user_id"`  // 0 = broadcast to all
	Group   string `gorm:"size:32;not null;index" json:"group"`   // system|order|marketing|im
	Title   string `gorm:"size:128;not null"  json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
	IsRead  int8   `gorm:"not null;default:0" json:"is_read"`
}
