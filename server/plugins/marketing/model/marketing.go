package model

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/model"
)

// ─── Coupon ──────────────────────────────────────────────────────────

const (
	CouponTypeFullReduce int8 = 1 // 满减
	CouponTypeDiscount   int8 = 2 // 折扣
	CouponTypeFree       int8 = 3 // 无门槛
)

type Coupon struct {
	model.Base
	Name            string     `gorm:"size:64;not null"              json:"name"`
	Type            int8       `gorm:"not null"                      json:"type"`
	MinAmount       float64    `gorm:"type:decimal(10,2);default:0"  json:"min_amount"`
	Discount        float64    `gorm:"type:decimal(10,2);not null"   json:"discount"`
	TotalCount      int        `gorm:"not null;default:0"            json:"total_count"`
	PerLimit        int        `gorm:"not null;default:1"            json:"per_limit"`
	Stackable       bool       `gorm:"not null;default:false"        json:"stackable"`        // can stack with other coupons?
	Scope           string     `gorm:"size:32;default:'all'"         json:"scope"`             // all|category|product
	ScopeIDs        string     `gorm:"type:json"                     json:"scope_ids"`         // JSON array of IDs
	ExcludeActivity bool       `gorm:"not null;default:false"        json:"exclude_activity"`  // skip items already in an activity?
	StartAt         *time.Time `json:"start_at"`
	EndAt           *time.Time `json:"end_at"`
	Status          int8       `gorm:"not null;default:1"            json:"status"`
}

type CouponUser struct {
	model.Base
	CouponID uint64     `gorm:"not null;index" json:"coupon_id"`
	UserID   uint64     `gorm:"not null;index" json:"user_id"`
	Status   int8       `gorm:"not null"       json:"status"` // 1未使用 2已使用 3已过期
	UsedAt   *time.Time `json:"used_at"`
	OrderID  uint64     `gorm:"default:0"      json:"order_id"`
}

// ─── Activity ────────────────────────────────────────────────────────

const (
	ActivityTypeSeckill   = "seckill"    // 秒杀
	ActivityTypeGroupBuy  = "group_buy"  // 拼团
	ActivityTypeBargain   = "bargain"    // 砍价
	ActivityTypeFullSave  = "full_save"  // 满减
	ActivityTypeCustom    = "custom"     // 自定义活动（520专场等）
)

// PriceRule defines how an activity modifies the price.
const (
	PriceRuleFixedPrice   = "fixed_price"   // replace unit price
	PriceRuleDiscountRate = "discount_rate"  // multiply (e.g. 0.8 = 80%)
	PriceRuleReduce       = "reduce"         // subtract fixed amount
)

type Activity struct {
	model.Base
	Type         string          `gorm:"size:32;not null;index"         json:"type"`
	Name         string          `gorm:"size:64;not null"               json:"name"`
	PriceRule    string          `gorm:"size:32"                        json:"price_rule"`   // fixed_price|discount_rate|reduce
	PriceValue   float64         `gorm:"type:decimal(10,2)"             json:"price_value"`  // the value for PriceRule
	ProductScope string          `gorm:"size:32;default:'selected'"     json:"product_scope"` // all|selected|category
	Exclusive    bool            `gorm:"not null;default:false"         json:"exclusive"`     // if true, no other discounts apply
	Config       json.RawMessage `gorm:"type:json"                      json:"config"`        // type-specific config (see below)
	StartAt      *time.Time      `json:"start_at"`
	EndAt        *time.Time      `json:"end_at"`
	Status       int8            `gorm:"not null;default:1"             json:"status"`
}

/*
Config per type:
  seckill:   {} (simple)
  group_buy: {"group_size": 3, "expire_hours": 24}
  bargain:   {"floor_price": 0.01, "max_helpers": 20, "expire_hours": 24}
  full_save: {"rules": [{"min": 100, "reduce": 20}, {"min": 200, "reduce": 50}]}
  custom:    {"label": "520专场", "banner": "https://..."}
*/

type ActivityProduct struct {
	model.Base
	ActivityID    uint64  `gorm:"not null;index"             json:"activity_id"`
	ProductID     uint64  `gorm:"not null"                   json:"product_id"`
	SkuID         uint64  `gorm:"default:0"                  json:"sku_id"`
	ActivityPrice float64 `gorm:"type:decimal(10,2)"         json:"activity_price"`
	ActivityStock int     `gorm:"default:0"                  json:"activity_stock"`
}

// ─── Group Buy ───────────────────────────────────────────────────────

type GroupBuyOrder struct {
	model.Base
	ActivityID  uint64     `gorm:"not null;index"  json:"activity_id"`
	LeaderID    uint64     `gorm:"not null"        json:"leader_id"`     // user who started the group
	GroupSize   int        `gorm:"not null"        json:"group_size"`    // required members
	JoinedCount int        `gorm:"not null;default:1" json:"joined_count"`
	Status      int8       `gorm:"not null"        json:"status"`        // 1=拼团中 2=成功 3=失败
	ExpireAt    *time.Time `json:"expire_at"`
}

type GroupBuyMember struct {
	model.Base
	GroupOrderID uint64 `gorm:"not null;index" json:"group_order_id"`
	UserID       uint64 `gorm:"not null"       json:"user_id"`
	OrderID      uint64 `gorm:"default:0"      json:"order_id"`
}

// ─── Bargain (砍价) ──────────────────────────────────────────────────

type BargainOrder struct {
	model.Base
	ActivityID   uint64  `gorm:"not null;index"  json:"activity_id"`
	UserID       uint64  `gorm:"not null;index"  json:"user_id"`
	ProductID    uint64  `gorm:"not null"        json:"product_id"`
	OriginalPrice float64 `gorm:"type:decimal(10,2);not null" json:"original_price"`
	CurrentPrice float64 `gorm:"type:decimal(10,2);not null" json:"current_price"`
	FloorPrice   float64 `gorm:"type:decimal(10,2);not null" json:"floor_price"`
	HelperCount  int     `gorm:"not null;default:0" json:"helper_count"`
	MaxHelpers   int     `gorm:"not null"        json:"max_helpers"`
	Status       int8    `gorm:"not null"        json:"status"` // 1=砍价中 2=已完成 3=已过期
	ExpireAt     *time.Time `json:"expire_at"`
}

type BargainHelper struct {
	model.Base
	BargainOrderID uint64  `gorm:"not null;index" json:"bargain_order_id"`
	UserID         uint64  `gorm:"not null"       json:"user_id"`
	CutAmount      float64 `gorm:"type:decimal(10,2);not null" json:"cut_amount"`
}

// ─── Distribution (分销) ─────────────────────────────────────────────

type Distributor struct {
	model.Base
	UserID    uint64  `gorm:"not null;uniqueIndex" json:"user_id"`
	ParentID  uint64  `gorm:"not null;default:0;index" json:"parent_id"`  // 上级分销商
	Level     int     `gorm:"not null;default:1" json:"level"`            // 1=一级 2=二级
	TotalEarn float64 `gorm:"type:decimal(10,2);default:0" json:"total_earn"`
	Balance   float64 `gorm:"type:decimal(10,2);default:0" json:"balance"`  // 可提现余额
	Status    int8    `gorm:"not null;default:1" json:"status"`
}

type DistributionCommission struct {
	model.Base
	OrderID       uint64  `gorm:"not null;index" json:"order_id"`
	DistributorID uint64  `gorm:"not null;index" json:"distributor_id"`
	Level         int     `gorm:"not null"       json:"level"`   // 1=直推 2=间推
	Amount        float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        int8    `gorm:"not null"       json:"status"`  // 1=待结算 2=已结算 3=已退回
}

type DistributionConfig struct {
	Level1Rate float64 `json:"level1_rate"` // e.g. 0.10 = 10%
	Level2Rate float64 `json:"level2_rate"` // e.g. 0.05 = 5%
}

// ─── Points ──────────────────────────────────────────────────────────

type PointsLog struct {
	model.Base
	UserID uint64 `gorm:"not null;index" json:"user_id"`
	Type   int8   `gorm:"not null"       json:"type"` // 1消费获取 2兑换消耗 3管理员调整
	Points int    `gorm:"not null"       json:"points"`
	Remark string `gorm:"size:128"       json:"remark"`
}
