package model

// User is the end-customer account.
type User struct {
	Base
	MerchantID uint64 `gorm:"not null;default:0;index" json:"merchant_id"` // reserved for multi-tenant
	Phone      string `gorm:"size:20;uniqueIndex"      json:"phone"`
	Nickname   string `gorm:"size:64"                  json:"nickname"`
	Avatar     string `gorm:"size:500"                 json:"avatar"`
	Points     int    `gorm:"not null;default:0"       json:"points"`
	Status     int8   `gorm:"not null;default:1"       json:"status"` // 1=active 0=banned
}
