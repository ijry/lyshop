package model

import (
	"time"

	"github.com/ijry/lyshop/model"
)

const (
	WarehouseStatusDisabled int8 = 0
	WarehouseStatusEnabled  int8 = 1

	DocTypeInbound  = "inbound"
	DocTypeOutbound = "outbound"

	DocStatusDraft     = "draft"
	DocStatusCompleted = "completed"
	DocStatusCanceled  = "canceled"
)

func IsValidDocType(v string) bool {
	return v == DocTypeInbound || v == DocTypeOutbound
}

func IsValidDocStatus(v string) bool {
	return v == DocStatusDraft || v == DocStatusCompleted || v == DocStatusCanceled
}

type Warehouse struct {
	model.Base
	Code    string `gorm:"size:32;not null;uniqueIndex" json:"code"`
	Name    string `gorm:"size:64;not null"             json:"name"`
	Address string `gorm:"size:255"                     json:"address"`
	Contact string `gorm:"size:64"                      json:"contact"`
	Phone   string `gorm:"size:20"                      json:"phone"`
	Status  int8   `gorm:"not null;default:1;index"     json:"status"`
}

func (Warehouse) TableName() string { return "warehouse" }

type InventoryStock struct {
	model.Base
	WarehouseID uint64 `gorm:"not null;index;uniqueIndex:uk_inventory_warehouse_sku" json:"warehouse_id"`
	SkuID       uint64 `gorm:"not null;index;uniqueIndex:uk_inventory_warehouse_sku" json:"sku_id"`
	Qty         int    `gorm:"not null;default:0"                                      json:"qty"`
	ReservedQty int    `gorm:"not null;default:0"                                      json:"reserved_qty"`
	SafeQty     int    `gorm:"not null;default:0"                                      json:"safe_qty"`
}

func (InventoryStock) TableName() string { return "inventory_stock" }

type InventoryMovement struct {
	model.Base
	WarehouseID uint64    `gorm:"not null;index"          json:"warehouse_id"`
	SkuID       uint64    `gorm:"not null;index"          json:"sku_id"`
	DocID       uint64    `gorm:"not null;index"          json:"doc_id"`
	DocNo       string    `gorm:"size:64;not null;index"  json:"doc_no"`
	BizType     string    `gorm:"size:32;not null;index"  json:"biz_type"`
	ChangeQty   int       `gorm:"not null"                json:"change_qty"`
	BeforeQty   int       `gorm:"not null"                json:"before_qty"`
	AfterQty    int       `gorm:"not null"                json:"after_qty"`
	OccurredAt  time.Time `gorm:"not null;index"          json:"occurred_at"`
	Remark      string    `gorm:"size:255"                json:"remark"`
}

func (InventoryMovement) TableName() string { return "inventory_movement" }

type InventoryDoc struct {
	model.Base
	DocNo       string     `gorm:"size:64;not null;uniqueIndex" json:"doc_no"`
	DocType     string     `gorm:"size:32;not null;index"       json:"doc_type"`
	Status      string     `gorm:"size:32;not null;index"       json:"status"`
	WarehouseID uint64     `gorm:"not null;index"               json:"warehouse_id"`
	Remark      string     `gorm:"size:255"                     json:"remark"`
	CompletedAt *time.Time `json:"completed_at"`
	CanceledAt  *time.Time `json:"canceled_at"`
}

func (InventoryDoc) TableName() string { return "inventory_doc" }

type InventoryDocItem struct {
	model.Base
	DocID  uint64 `gorm:"not null;index" json:"doc_id"`
	SkuID  uint64 `gorm:"not null;index" json:"sku_id"`
	Qty    int    `gorm:"not null"       json:"qty"`
	Remark string `gorm:"size:255"       json:"remark"`
}

func (InventoryDocItem) TableName() string { return "inventory_doc_item" }
