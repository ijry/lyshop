package model

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/model"
)

// CouponType constants
const (
	CouponTypeFullReduce int8 = 1 // 满减
	CouponTypeDiscount   int8 = 2 // 折扣
	CouponTypeFree       int8 = 3 // 无门槛
)

// Coupon is a coupon template.
type Coupon struct {
	model.Base
	Name       string     `gorm:"size:64;not null"            json:"name"`
	Type       int8       `gorm:"not null"                    json:"type"`
	MinAmount  float64    `gorm:"type:decimal(10,2);default:0" json:"min_amount"`
	Discount   float64    `gorm:"type:decimal(10,2);not null"  json:"discount"`
	TotalCount int        `gorm:"not null;default:0"           json:"total_count"` // 0=unlimited
	PerLimit   int        `gorm:"not null;default:1"           json:"per_limit"`
	StartAt    *time.Time `json:"start_at"`
	EndAt      *time.Time `json:"end_at"`
	Status     int8       `gorm:"not null;default:1"           json:"status"`
}

// CouponUser is a user's coupon (claimed instance).
type CouponUser struct {
	model.Base
	CouponID uint64     `gorm:"not null;index"  json:"coupon_id"`
	UserID   uint64     `gorm:"not null;index"  json:"user_id"`
	Status   int8       `gorm:"not null"        json:"status"` // 1未使用 2已使用 3已过期
	UsedAt   *time.Time `json:"used_at"`
	OrderID  uint64     `gorm:"default:0"       json:"order_id"`
}

// ActivityType constants
const (
	ActivityTypeSeckill  int8 = 1
	ActivityTypeFullSave int8 = 2
)

// Activity is a marketing activity (seckill or full-reduce).
type Activity struct {
	model.Base
	Type    int8            `gorm:"not null"                  json:"type"`
	Name    string          `gorm:"size:64;not null"          json:"name"`
	Config  json.RawMessage `gorm:"type:json"                 json:"config"`
	StartAt *time.Time      `json:"start_at"`
	EndAt   *time.Time      `json:"end_at"`
	Status  int8            `gorm:"not null;default:1"        json:"status"`
}

// ActivityProduct links a product/SKU to an activity with special pricing.
type ActivityProduct struct {
	model.Base
	ActivityID    uint64  `gorm:"not null;index"               json:"activity_id"`
	ProductID     uint64  `gorm:"not null"                     json:"product_id"`
	SkuID         uint64  `gorm:"default:0"                    json:"sku_id"`
	ActivityPrice float64 `gorm:"type:decimal(10,2)"           json:"activity_price"`
	ActivityStock int     `gorm:"default:0"                    json:"activity_stock"`
}

// PointsLog records point changes for a user.
type PointsLog struct {
	model.Base
	UserID uint64 `gorm:"not null;index"  json:"user_id"`
	Type   int8   `gorm:"not null"        json:"type"` // 1消费获取 2兑换消耗 3管理员调整
	Points int    `gorm:"not null"        json:"points"`
	Remark string `gorm:"size:128"        json:"remark"`
}
