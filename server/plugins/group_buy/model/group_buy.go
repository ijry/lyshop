package model

import (
	"time"

	"github.com/ijry/lyshop/core/model"
)

// GroupBuyActivity 拼团活动
type GroupBuyActivity struct {
	model.Base
	Name        string     `gorm:"size:128;not null" json:"name"`
	GroupSize   int        `gorm:"not null;default:2" json:"group_size"`     // 成团人数
	ExpireHours int        `gorm:"not null;default:24" json:"expire_hours"`  // 过期小时数
	StartAt     *time.Time `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	Status      int8       `gorm:"not null;default:1;index" json:"status"`   // 1=启用 0=禁用
	Sort        int        `gorm:"not null;default:0" json:"sort"`
}

func (GroupBuyActivity) TableName() string {
	return "group_buy_activities"
}

// GroupBuyProduct 拼团商品
type GroupBuyProduct struct {
	model.Base
	ActivityID      uint64  `gorm:"not null;index" json:"activity_id"`
	ProductID       uint64  `gorm:"not null;index" json:"product_id"`
	SkuID           uint64  `gorm:"default:0;index" json:"sku_id"`
	GroupPrice      float64 `gorm:"type:decimal(10,2);not null" json:"group_price"`
	LimitPerOrder   int     `gorm:"not null;default:0" json:"limit_per_order"`
	TotalStockLimit int     `gorm:"not null;default:0" json:"total_stock_limit"`
	SoldQty         int     `gorm:"not null;default:0" json:"sold_qty"`
	Sort            int     `gorm:"not null;default:0" json:"sort"`
}

func (GroupBuyProduct) TableName() string {
	return "group_buy_products"
}

// GroupBuyOrder 拼团订单
type GroupBuyOrder struct {
	model.Base
	ActivityID  uint64     `gorm:"not null;index" json:"activity_id"`
	ProductID   uint64     `gorm:"not null;index" json:"product_id"`
	SkuID       uint64     `gorm:"not null;index" json:"sku_id"`
	LeaderID    uint64     `gorm:"not null;index" json:"leader_id"`      // 团长用户ID
	GroupSize   int        `gorm:"not null" json:"group_size"`           // 需要人数
	JoinedCount int        `gorm:"not null;default:1" json:"joined_count"` // 已参团人数
	Status      string     `gorm:"size:32;not null;default:'pending';index" json:"status"` // pending|success|failed
	ExpireAt    *time.Time `json:"expire_at"`
}

func (GroupBuyOrder) TableName() string {
	return "group_buy_orders"
}

// GroupBuyMember 拼团成员
type GroupBuyMember struct {
	model.Base
	GroupOrderID uint64 `gorm:"not null;index" json:"group_order_id"`
	UserID       uint64 `gorm:"not null;index" json:"user_id"`
	OrderID      uint64 `gorm:"not null;index" json:"order_id"` // 关联的订单ID
	IsLeader     bool   `gorm:"not null;default:false" json:"is_leader"`
}

func (GroupBuyMember) TableName() string {
	return "group_buy_members"
}
