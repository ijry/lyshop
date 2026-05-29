package model

import (
	"time"

	"github.com/ijry/lyshop/model"
)

const (
	ReservationStatusReserved  = "reserved"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusReleased  = "released"
)

type InventoryReservation struct {
	model.Base
	BizType     string     `gorm:"size:32;not null;uniqueIndex:uk_inventory_reservation_biz_sku,priority:1;index" json:"biz_type"`
	BizNo       string     `gorm:"size:64;not null;uniqueIndex:uk_inventory_reservation_biz_sku,priority:2;index" json:"biz_no"`
	WarehouseID uint64     `gorm:"not null;index"                                                               json:"warehouse_id"`
	SkuID       uint64     `gorm:"not null;uniqueIndex:uk_inventory_reservation_biz_sku,priority:3;index"      json:"sku_id"`
	Qty         int        `gorm:"not null"                                                                     json:"qty"`
	Status      string     `gorm:"size:16;not null;index"                                                      json:"status"`
	ExpiredAt   *time.Time `json:"expired_at"`
	Remark      string     `gorm:"size:255"                                                                     json:"remark"`
}

func (InventoryReservation) TableName() string { return "inventory_reservation" }
