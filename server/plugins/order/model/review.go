package model

import "github.com/ijry/lyshop/model"

// OrderReview is the root review of one order item.
type OrderReview struct {
	model.Base
	OrderID        uint64 `gorm:"not null;index"               json:"order_id"`
	OrderItemID    uint64 `gorm:"not null;uniqueIndex"         json:"order_item_id"`
	ProductID      uint64 `gorm:"not null;index"               json:"product_id"`
	UserID         uint64 `gorm:"not null;index"               json:"user_id"`
	MerchantID     uint64 `gorm:"not null;default:0;index"     json:"merchant_id"`
	ProductScore   int8   `gorm:"not null;default:5"           json:"product_score"`
	LogisticsScore int8   `gorm:"not null;default:5"           json:"logistics_score"`
	Content        string `gorm:"type:text"                    json:"content"`
	EditedTimes    int    `gorm:"not null;default:0"           json:"edited_times"`
}

// OrderReviewAppend is an append-only follow-up review under a root review.
type OrderReviewAppend struct {
	model.Base
	ReviewID uint64 `gorm:"not null;index"               json:"review_id"`
	UserID   uint64 `gorm:"not null;index"               json:"user_id"`
	Content  string `gorm:"type:text"                    json:"content"`
}

// OrderReviewReply is a single admin reply for one root review.
type OrderReviewReply struct {
	model.Base
	ReviewID uint64 `gorm:"not null;uniqueIndex"         json:"review_id"`
	AdminID  uint64 `gorm:"not null;index"               json:"admin_id"`
	Content  string `gorm:"type:text"                    json:"content"`
}
