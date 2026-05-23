package model

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/model"
)

// OrderStatus constants
const (
	OrderStatusPending   int8 = 1 // 待付款
	OrderStatusPaid      int8 = 2 // 待发货
	OrderStatusShipped   int8 = 3 // 待收货
	OrderStatusCompleted int8 = 4 // 已完成
	OrderStatusAfterSale int8 = 5 // 售后
)

type Order struct {
	model.Base
	OrderNo        string          `gorm:"size:64;uniqueIndex;not null" json:"order_no"`
	UserID         uint64          `gorm:"not null;index"               json:"user_id"`
	MerchantID     uint64          `gorm:"not null;default:0"           json:"merchant_id"`
	Status         int8            `gorm:"not null;index"               json:"status"`
	PaymentMethod  string          `gorm:"size:32"                      json:"payment_method"`
	GoodsAmount    float64         `gorm:"type:decimal(10,2);not null"  json:"goods_amount"`
	DiscountAmount float64         `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	FreightAmount  float64         `gorm:"type:decimal(10,2);default:0" json:"freight_amount"`
	TotalAmount    float64         `gorm:"type:decimal(10,2);not null"  json:"total_amount"`
	AddressSnapshot json.RawMessage `gorm:"type:json"                   json:"address_snapshot"`
	Remark         string          `gorm:"size:255"                     json:"remark"`
	TrackingNo     string          `gorm:"size:128"                     json:"tracking_no"`
	PaidAt         *time.Time      `json:"paid_at"`
}

type OrderItem struct {
	model.Base
	OrderID   uint64          `gorm:"not null;index"               json:"order_id"`
	ProductID uint64          `gorm:"not null"                     json:"product_id"`
	SkuID     uint64          `gorm:"not null"                     json:"sku_id"`
	Title     string          `gorm:"size:255;not null"            json:"title"`
	Cover     string          `gorm:"size:500"                     json:"cover"`
	Attrs     json.RawMessage `gorm:"type:json"                    json:"attrs"`
	Price     float64         `gorm:"type:decimal(10,2);not null"  json:"price"`
	Qty       int             `gorm:"not null"                     json:"qty"`
}

type OrderPayment struct {
	model.Base
	OrderID    uint64     `gorm:"not null;index"               json:"order_id"`
	Driver     string     `gorm:"size:32;not null"             json:"driver"`
	TradeNo    string     `gorm:"size:128"                     json:"trade_no"`
	Amount     float64    `gorm:"type:decimal(10,2);not null"  json:"amount"`
	Status     int8       `gorm:"not null"                     json:"status"` // 1待支付 2已支付 3已退款
	NotifiedAt *time.Time `json:"notified_at"`
}

type OrderRefund struct {
	model.Base
	OrderID  uint64  `gorm:"not null;index"              json:"order_id"`
	Reason   string  `gorm:"size:255"                    json:"reason"`
	Amount   float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status   int8    `gorm:"not null"                    json:"status"` // 1申请中 2已退款 3已拒绝
	RefundNo string  `gorm:"size:128"                    json:"refund_no"`
}
