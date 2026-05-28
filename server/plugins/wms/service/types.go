package service

import (
	"fmt"
	"time"
)

type WarehouseListQuery struct {
	Page    int
	Size    int
	Keyword string
	Status  *int8
}

type StockListQuery struct {
	Page        int
	Size        int
	WarehouseID uint64
	SkuID       uint64
	WarningOnly bool
}

type DocListQuery struct {
	Page        int
	Size        int
	WarehouseID uint64
	DocType     string
	Status      string
	DocNo       string
	StartAt     *time.Time
	EndAt       *time.Time
}

type MovementListQuery struct {
	Page        int
	Size        int
	WarehouseID uint64
	SkuID       uint64
	BizType     string
	DocNo       string
	StartAt     *time.Time
	EndAt       *time.Time
}

type DocItemInput struct {
	SkuID  uint64 `json:"sku_id"`
	Qty    int    `json:"qty"`
	Remark string `json:"remark"`
}

type CreateDocInput struct {
	WarehouseID uint64         `json:"warehouse_id"`
	DocType     string         `json:"doc_type"`
	Remark      string         `json:"remark"`
	Items       []DocItemInput `json:"items"`
}

type UpdateDocInput struct {
	WarehouseID uint64         `json:"warehouse_id"`
	DocType     string         `json:"doc_type"`
	Remark      string         `json:"remark"`
	Items       []DocItemInput `json:"items"`
}

func normalizePage(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	return page, size
}

func genDocNo(docType string) string {
	prefix := "WD"
	if docType == "inbound" {
		prefix = "WI"
	} else if docType == "outbound" {
		prefix = "WO"
	}
	return fmt.Sprintf("%s%s", prefix, time.Now().Format("20060102150405.000000"))
}
