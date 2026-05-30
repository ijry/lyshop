package model

import (
	"encoding/json"
	"time"

	"github.com/ijry/lyshop/server/model"
)

// PointsLog 积分日志
type PointsLog struct {
	model.Base
	UserID    uint64 `gorm:"not null;index" json:"user_id"`
	Type      int8   `gorm:"not null;index" json:"type"` // 1=签到 2=订单抵扣 3=兑换消耗 4=订单完成 5=管理员调整 6=过期扣除 7=活动奖励
	Points    int    `gorm:"not null" json:"points"`     // 正数=增加，负数=减少
	RelatedID uint64 `gorm:"default:0;index" json:"related_id"` // 关联ID（订单ID、兑换ID等）
	Remark    string `gorm:"size:255" json:"remark"`
}

// PointsProduct 积分商品
type PointsProduct struct {
	model.Base
	Title       string `gorm:"size:128;not null" json:"title"`
	Type        string `gorm:"size:32;not null;index" json:"type"` // coupon=优惠券 physical=实物 virtual=虚拟
	PointsPrice int    `gorm:"not null" json:"points_price"`       // 积分价格
	Stock       int    `gorm:"not null;default:0" json:"stock"`    // 库存（0=无限）
	SoldCount   int    `gorm:"not null;default:0" json:"sold_count"` // 已兑换数量
	Cover       string `gorm:"size:500" json:"cover"`
	Images      string `gorm:"type:text" json:"images"` // JSON数组
	Description string `gorm:"type:text" json:"description"`
	Status      int8   `gorm:"not null;default:1;index" json:"status"` // 1=上架 0=下架
	Sort        int    `gorm:"not null;default:0" json:"sort"`

	// 兑换限制
	LimitPerUser int `gorm:"not null;default:0" json:"limit_per_user"` // 每人限兑（0=不限）
	LimitPerDay  int `gorm:"not null;default:0" json:"limit_per_day"`  // 每日限兑（0=不限）

	// 优惠券类型专用字段
	CouponID uint64 `gorm:"default:0;index" json:"coupon_id"` // 关联的优惠券ID

	// 实物类型专用字段
	NeedAddress bool `gorm:"not null;default:false" json:"need_address"` // 是否需要收货地址

	// 虚拟类型专用字段
	VirtualContent string `gorm:"type:text" json:"virtual_content"` // 虚拟商品内容（兑换码等）
}

// PointsExchange 兑换记录
type PointsExchange struct {
	model.Base
	ExchangeNo      string          `gorm:"size:64;uniqueIndex;not null" json:"exchange_no"` // 兑换单号
	UserID          uint64          `gorm:"not null;index" json:"user_id"`
	ProductID       uint64          `gorm:"not null;index" json:"product_id"`
	ProductTitle    string          `gorm:"size:128;not null" json:"product_title"` // 冗余，防止商品删除
	ProductType     string          `gorm:"size:32;not null" json:"product_type"`
	ProductCover    string          `gorm:"size:500" json:"product_cover"` // 冗余商品封面
	PointsCost      int             `gorm:"not null" json:"points_cost"`
	Qty             int             `gorm:"not null;default:1" json:"qty"`
	Status          string          `gorm:"size:32;not null;index" json:"status"` // pending_ship=待发货 shipped=已发货 completed=已完成 canceled=已取消

	// 收货地址（实物商品）
	AddressSnapshot json.RawMessage `gorm:"type:json" json:"address_snapshot"`
	TrackingNo      string          `gorm:"size:128" json:"tracking_no"` // 物流单号

	// 优惠券（优惠券类商品）
	CouponUserID uint64 `gorm:"default:0;index" json:"coupon_user_id"` // 发放的优惠券记录ID

	// 虚拟商品
	VirtualContent string `gorm:"type:text" json:"virtual_content"` // 虚拟商品内容

	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Remark      string     `gorm:"size:255" json:"remark"`
}

// PointsConfig 积分规则配置
type PointsConfig struct {
	model.Base
	Key         string `gorm:"size:64;uniqueIndex;not null" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Description string `gorm:"size:255" json:"description"`
}

// TableName 指定表名
func (PointsLog) TableName() string {
	return "points_logs"
}

func (PointsProduct) TableName() string {
	return "points_products"
}

func (PointsExchange) TableName() string {
	return "points_exchanges"
}

func (PointsConfig) TableName() string {
	return "points_configs"
}
