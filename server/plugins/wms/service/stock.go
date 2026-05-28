package service

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/core/db"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
)

type StockView struct {
	wmsmodel.InventoryStock
	IsWarning  bool `json:"is_warning"`
	WarningGap int  `json:"warning_gap"`
}

func ListStocks(ctx context.Context, q StockListQuery) ([]StockView, int64, error) {
	page, size := normalizePage(q.Page, q.Size)
	tx := db.DB.WithContext(ctx).Model(&wmsmodel.InventoryStock{})
	if q.WarehouseID > 0 {
		tx = tx.Where("warehouse_id = ?", q.WarehouseID)
	}
	if q.SkuID > 0 {
		tx = tx.Where("sku_id = ?", q.SkuID)
	}
	if q.WarningOnly {
		tx = tx.Where("qty < safe_qty")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []wmsmodel.InventoryStock
	if err := tx.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	list := make([]StockView, 0, len(rows))
	for _, row := range rows {
		gap := row.SafeQty - row.Qty
		if gap < 0 {
			gap = 0
		}
		list = append(list, StockView{
			InventoryStock: row,
			IsWarning:      row.Qty < row.SafeQty,
			WarningGap:     gap,
		})
	}
	return list, total, nil
}

func UpdateStockSafety(ctx context.Context, id uint64, safeQty int) error {
	if id == 0 {
		return fmt.Errorf("库存ID不能为空")
	}
	if safeQty < 0 {
		return fmt.Errorf("安全库存不能小于0")
	}
	tx := db.DB.WithContext(ctx).Model(&wmsmodel.InventoryStock{}).Where("id = ?", id).Update("safe_qty", safeQty)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("库存记录不存在")
	}
	return nil
}
