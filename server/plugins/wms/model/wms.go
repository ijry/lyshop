package model

import "github.com/ijry/lyshop/model"

type Warehouse struct {
	model.Base
	Name    string `gorm:"size:64;not null" json:"name"`
	Address string `gorm:"size:255"         json:"address"`
	Contact string `gorm:"size:64"          json:"contact"`
	Phone   string `gorm:"size:20"          json:"phone"`
	Status  int8   `gorm:"not null;default:1" json:"status"`
}

type WmsStock struct {
	model.Base
	WarehouseID uint64 `gorm:"not null;uniqueIndex:uk_warehouse_sku" json:"warehouse_id"`
	SkuID       uint64 `gorm:"not null;uniqueIndex:uk_warehouse_sku" json:"sku_id"`
	Qty         int    `gorm:"not null;default:0"                   json:"qty"`
	SafeQty     int    `gorm:"not null;default:0"                   json:"safe_qty"`
}

// WmsInbound is an inbound order header.
type WmsInbound struct {
	model.Base
	WarehouseID uint64 `gorm:"not null;index" json:"warehouse_id"`
	Type        int8   `gorm:"not null"       json:"type"`   // 1=采购 2=退货
	Status      int8   `gorm:"not null"       json:"status"` // 1=待入库 2=已完成
	Remark      string `gorm:"size:255"       json:"remark"`
}

type WmsInboundItem struct {
	model.Base
	InboundID uint64 `gorm:"not null;index" json:"inbound_id"`
	SkuID     uint64 `gorm:"not null"       json:"sku_id"`
	Qty       int    `gorm:"not null"       json:"qty"`
}

// WmsOutbound is an outbound order header.
type WmsOutbound struct {
	model.Base
	WarehouseID uint64 `gorm:"not null;index" json:"warehouse_id"`
	Type        int8   `gorm:"not null"       json:"type"`   // 1=订单出库 2=调拨出库
	RefID       uint64 `gorm:"default:0"      json:"ref_id"` // order_id or transfer_id
	Status      int8   `gorm:"not null"       json:"status"`
}

type WmsOutboundItem struct {
	model.Base
	OutboundID uint64 `gorm:"not null;index" json:"outbound_id"`
	SkuID      uint64 `gorm:"not null"       json:"sku_id"`
	Qty        int    `gorm:"not null"       json:"qty"`
}

// WmsStockLog records every stock change.
type WmsStockLog struct {
	model.Base
	WarehouseID uint64 `gorm:"not null;index"  json:"warehouse_id"`
	SkuID       uint64 `gorm:"not null;index"  json:"sku_id"`
	Type        string `gorm:"size:32;not null" json:"type"` // inbound/outbound/adjust
	Qty         int    `gorm:"not null"        json:"qty"`        // positive=in, negative=out
	BeforeQty   int    `gorm:"not null"        json:"before_qty"`
	AfterQty    int    `gorm:"not null"        json:"after_qty"`
	RefID       uint64 `gorm:"default:0"       json:"ref_id"`
	RefType     string `gorm:"size:32"         json:"ref_type"`
}
