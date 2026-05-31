package model

import (
	"time"

	"github.com/ijry/lyshop/core/model"
)

// BargainActivity 砍价活动
type BargainActivity struct {
	model.Base
	Name    string     `gorm:"size:128;not null" json:"name"`
	StartAt *time.Time `json:"start_at"`
	EndAt   *time.Time `json:"end_at"`
	Status  int8       `gorm:"not null;default:1;index" json:"status"` // 1=启用 0=禁用
	Sort    int        `gorm:"not null;default:0" json:"sort"`
}

func (BargainActivity) TableName() string {
	return "bargain_activities"
}

// BargainProduct 砍价商品
type BargainProduct struct {
	model.Base
	ActivityID      uint64  `gorm:"not null;index" json:"activity_id"`
	ProductID       uint64  `gorm:"not null;index" json:"product_id"`
	SkuID           uint64  `gorm:"default:0;index" json:"sku_id"`
	StartPrice      float64 `gorm:"type:decimal(10,2);not null" json:"start_price"`       // 起始价格
	FloorPrice      float64 `gorm:"type:decimal(10,2);not null" json:"floor_price"`       // 底价
	MinCutAmount    float64 `gorm:"type:decimal(10,2);not null;default:0.01" json:"min_cut_amount"` // 最小砍价金额
	MaxCutAmount    float64 `gorm:"type:decimal(10,2);not null;default:10" json:"max_cut_amount"`   // 最大砍价金额
	MaxHelpers      int     `gorm:"not null;default:50" json:"max_helpers"`               // 最多助力人数
	ExpireHours     int     `gorm:"not null;default:24" json:"expire_hours"`              // 过期小时数
	TotalStockLimit int     `gorm:"not null;default:0" json:"total_stock_limit"`
	SoldQty         int     `gorm:"not null;default:0" json:"sold_qty"`
	Sort            int     `gorm:"not null;default:0" json:"sort"`
}

func (BargainProduct) TableName() string {
	return "bargain_products"
}

// BargainOrder 砍价订单
type BargainOrder struct {
	model.Base
	ActivityID   uint64     `gorm:"not null;index" json:"activity_id"`
	ProductID    uint64     `gorm:"not null;index" json:"product_id"`
	SkuID        uint64     `gorm:"not null;index" json:"sku_id"`
	UserID       uint64     `gorm:"not null;index" json:"user_id"`           // 发起人
	StartPrice   float64    `gorm:"type:decimal(10,2);not null" json:"start_price"`
	FloorPrice   float64    `gorm:"type:decimal(10,2);not null" json:"floor_price"`
	CurrentPrice float64    `gorm:"type:decimal(10,2);not null" json:"current_price"` // 当前价格
	HelpCount    int        `gorm:"not null;default:0" json:"help_count"`             // 助力次数
	Status       string     `gorm:"size:32;not null;default:'pending';index" json:"status"` // pending|success|failed
	ExpireAt     *time.Time `json:"expire_at"`
}

func (BargainOrder) TableName() string {
	return "bargain_orders"
}

// BargainHelper 砍价助力记录
type BargainHelper struct {
	model.Base
	BargainOrderID uint64  `gorm:"not null;index" json:"bargain_order_id"`
	UserID         uint64  `gorm:"not null;index" json:"user_id"`
	CutAmount      float64 `gorm:"type:decimal(10,2);not null" json:"cut_amount"` // 砍掉的金额
}

func (BargainHelper) TableName() string {
	return "bargain_helpers"
}
