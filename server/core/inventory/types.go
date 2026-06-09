package inventory

import "time"

type ReserveItem struct {
	SkuID uint64
	Qty   int
}

type ReserveInput struct {
	BizType   string
	BizNo     string
	Items     []ReserveItem
	ExpiredAt *time.Time
}

type DeductInput struct {
	BizType string
	BizNo   string
	Items   []ReserveItem
}

type RestoreInput struct {
	BizType string
	BizNo   string
	Items   []ReserveItem
	Reason  string
}

type SyncSkuInput struct {
	SkuID  uint64
	Stock  int
	Source string
}

type SellableStock struct {
	SkuID         uint64 `json:"sku_id"`
	SellableStock int    `json:"sellable_stock"`
	Reserved      int    `json:"reserved"`
	OnHand        int    `json:"on_hand"`
}
