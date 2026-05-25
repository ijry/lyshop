package model

import (
	"time"

	"github.com/ijry/lyshop/model"
)

type Plan struct {
	model.Base
	Name           string  `gorm:"size:64;not null"              json:"name"`
	DurationMonths int     `gorm:"not null;default:12"           json:"duration_months"`
	Price          float64 `gorm:"type:decimal(10,2);not null"   json:"price"`
	RenewPrice     float64 `gorm:"type:decimal(10,2);default:0"  json:"renew_price"`
	Status         int8    `gorm:"not null;default:1"            json:"status"`
}

type Level struct {
	model.Base
	Name            string `gorm:"size:64;not null"           json:"name"`
	GrowthThreshold int64  `gorm:"not null;default:0;index"   json:"growth_threshold"`
	BenefitJSON     string `gorm:"type:json"                  json:"benefit_json"`
	Status          int8   `gorm:"not null;default:1"         json:"status"`
	Sort            int    `gorm:"not null;default:0"         json:"sort"`
}

type UserAsset struct {
	model.Base
	UserID         uint64     `gorm:"not null;uniqueIndex"     json:"user_id"`
	CurrentPlanID  uint64     `gorm:"not null;default:0"       json:"current_plan_id"`
	CurrentLevelID uint64     `gorm:"not null;default:0"       json:"current_level_id"`
	VipStartAt     *time.Time `json:"vip_start_at"`
	VipEndAt       *time.Time `json:"vip_end_at"`
	GrowthValue    int64      `gorm:"not null;default:0"       json:"growth_value"`
	Status         int8       `gorm:"not null;default:0"       json:"status"` // 0=非会员 1=有效
}

type GrowthLog struct {
	model.Base
	UserID         uint64 `gorm:"not null;index"               json:"user_id"`
	OrderID        uint64 `gorm:"not null;default:0;index"     json:"order_id"`
	EventType      string `gorm:"size:32;not null"             json:"event_type"` // order_paid|order_refund|admin_adjust
	GrowthDelta    int64  `gorm:"not null"                     json:"growth_delta"`
	BalanceAfter   int64  `gorm:"not null;default:0"           json:"balance_after"`
	IdempotencyKey string `gorm:"size:128;not null;uniqueIndex" json:"idempotency_key"`
	Remark         string `gorm:"size:128"                     json:"remark"`
}

type CouponRule struct {
	model.Base
	PlanID       uint64 `gorm:"not null;default:0;index"       json:"plan_id"`
	LevelID      uint64 `gorm:"not null;default:0;index"       json:"level_id"`
	CouponID     uint64 `gorm:"not null;index"                 json:"coupon_id"`
	MonthlyQuota int    `gorm:"not null;default:1"             json:"monthly_quota"`
	ClaimMode    string `gorm:"size:16;not null;default:'manual'" json:"claim_mode"` // manual
	Status       int8   `gorm:"not null;default:1"             json:"status"`
}

type CouponClaim struct {
	model.Base
	UserID        uint64     `gorm:"not null;uniqueIndex:ux_vip_claim_period" json:"user_id"`
	RuleID        uint64     `gorm:"not null;uniqueIndex:ux_vip_claim_period;index" json:"rule_id"`
	PeriodYYYYMM  string     `gorm:"size:6;not null;uniqueIndex:ux_vip_claim_period" json:"period_yyyymm"`
	ClaimedCount  int        `gorm:"not null;default:0"        json:"claimed_count"`
	LastClaimedAt *time.Time `json:"last_claimed_at"`
}

type SkuPrice struct {
	model.Base
	ProductID       uint64  `gorm:"not null;index"                            json:"product_id"`
	SkuID           uint64  `gorm:"not null;default:0;index"                  json:"sku_id"`
	LevelID         uint64  `gorm:"not null;index"                            json:"level_id"`
	VipPrice        float64 `gorm:"type:decimal(10,2);default:0"              json:"vip_price"`
	VipDiscountRate float64 `gorm:"type:decimal(10,4);default:0"              json:"vip_discount_rate"`
	Status          int8    `gorm:"not null;default:1"                        json:"status"`
}

type OrderBenefit struct {
	model.Base
	OrderID        uint64  `gorm:"not null;uniqueIndex"        json:"order_id"`
	UserID         uint64  `gorm:"not null;index"              json:"user_id"`
	VipDiscount    float64 `gorm:"type:decimal(10,2);default:0" json:"vip_discount"`
	GrowthGranted  int64   `gorm:"not null;default:0"          json:"growth_granted"`
	GrowthReverted int64   `gorm:"not null;default:0"          json:"growth_reverted"`
}
