package model

import (
	"time"

	"github.com/ijry/lyshop/model"
)

type ShipmentDirection string

const (
	ShipmentDirectionOutbound ShipmentDirection = "outbound"
	ShipmentDirectionInbound  ShipmentDirection = "inbound"
)

type ShipmentBizType string

const (
	ShipmentBizTypeInitial ShipmentBizType = "initial"
	ShipmentBizTypeReship  ShipmentBizType = "reship"
	ShipmentBizTypeReturn  ShipmentBizType = "return"
)

type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "pending"
	ShipmentStatusShipped   ShipmentStatus = "shipped"
	ShipmentStatusInTransit ShipmentStatus = "in_transit"
	ShipmentStatusSigned    ShipmentStatus = "signed"
	ShipmentStatusException ShipmentStatus = "exception"
)

type OrderShipment struct {
	model.Base
	OrderID         uint64     `gorm:"not null;index"                                json:"order_id"`
	AfterSaleCaseID uint64     `gorm:"index"                                         json:"after_sale_case_id"`
	Direction       string     `gorm:"size:16;not null;index"                        json:"direction"`
	BizType         string     `gorm:"size:16;not null;index"                        json:"biz_type"`
	Company         string     `gorm:"size:64"                                       json:"company"`
	TrackingNo      string     `gorm:"size:128;index"                                json:"tracking_no"`
	LogisticsStatus string     `gorm:"size:32;not null;default:'pending'"            json:"logistics_status"`
	Remark          string     `gorm:"size:255"                                      json:"remark"`
	ShippedAt       *time.Time `json:"shipped_at"`
	SignedAt        *time.Time `json:"signed_at"`
	CreatedByType   string     `gorm:"size:16;not null;default:'admin'"              json:"created_by_type"`
	CreatedByID     uint64     `gorm:"not null;default:0"                            json:"created_by_id"`
}
