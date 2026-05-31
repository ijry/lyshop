package model

import (
	"time"

	"github.com/ijry/lyshop/server/model"
)

// SeckillActivity 秒杀活动
type SeckillActivity struct {
	model.Base
	Name    string     `gorm:"size:128;not null" json:"name"`
	StartAt *time.Time `json:"start_at"`
	EndAt   *time.Time `json:"end_at"`
	Status  int8       `gorm:"not null;default:1;index" json:"status"` // 1=启用 0=禁用
	Sort    int        `gorm:"not null;default:0" json:"sort"`
}

// SeckillProduct 秒杀商品
type SeckillProduct struct {
	model.Base
	ActivityID      uint64  `gorm:"not null;index" json:"activity_id"`
	ProductID       uint64  `gorm:"not null;index" json:"product_id"`
	SkuID           uint64  `gorm:"default:0;index" json:"sku_id"` // 0表示全部SKU
	SeckillPrice    float64 `gorm:"type:decimal(10,2);not null" json:"seckill_price"`
	LimitPerOrder   int     `gorm:"not null;default:0" json:"limit_per_order"`   // 单笔限购（0=不限）
	TotalStockLimit int     `gorm:"not null;default:0" json:"total_stock_limit"` // 活动库存上限（0=不限）
	SoldQty         int     `gorm:"not null;default:0" json:"sold_qty"`          // 已售数量
	Sort            int     `gorm:"not null;default:0" json:"sort"`
}

// TableName 指定表名
func (SeckillActivity) TableName() string {
	return "seckill_activities"
}

func (SeckillProduct) TableName() string {
	return "seckill_products"
}
