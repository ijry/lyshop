package model

import "github.com/ijry/lyshop/model"

const (
	SessionStatusWaiting  int8 = 1
	SessionStatusOngoing  int8 = 2
	SessionStatusClosed   int8 = 3

	MsgTypeText        = "text"
	MsgTypeImage       = "image"
	MsgTypeProductCard = "product_card"
	MsgTypeOrderCard   = "order_card"

	SenderUser  int8 = 1
	SenderStaff int8 = 2
)

type ImSession struct {
	model.Base
	UserID      uint64 `gorm:"not null;index"     json:"user_id"`
	StaffID     uint64 `gorm:"not null;default:0;index" json:"staff_id"`
	Status      int8   `gorm:"not null"           json:"status"`
	LastMsg     string `gorm:"size:255"           json:"last_msg"`
	UnreadCount int    `gorm:"not null;default:0" json:"unread_count"`
}

type ImMessage struct {
	model.Base
	SessionID  uint64 `gorm:"not null;index"    json:"session_id"`
	SenderType int8   `gorm:"not null"          json:"sender_type"`
	SenderID   uint64 `gorm:"not null"          json:"sender_id"`
	Type       string `gorm:"size:32;not null"  json:"type"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Extra      string `gorm:"type:json"         json:"extra"`
	IsRead     int8   `gorm:"not null;default:0" json:"is_read"`
}

type ImAutoReply struct {
	model.Base
	Keyword   string `gorm:"size:128;not null" json:"keyword"`
	MatchType int8   `gorm:"not null"          json:"match_type"` // 1精确 2包含 3正则
	Reply     string `gorm:"type:text;not null" json:"reply"`
	Sort      int    `gorm:"not null;default:0" json:"sort"`
	Status    int8   `gorm:"not null;default:1" json:"status"`
}
