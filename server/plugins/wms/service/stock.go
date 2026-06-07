package service

import (
	"context"
	"errors"

	"github.com/ijry/lyshop/core/db"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	"gorm.io/gorm"
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
		return nil, 0, WrapDBError("查询库存总数失败", err)
	}

	var rows []wmsmodel.InventoryStock
	if err := tx.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, 0, WrapDBError("查询库存列表失败", err)
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

func ListStocksBySkuIDs(ctx context.Context, skuIDs []uint64) ([]wmsmodel.InventoryStock, error) {
	if len(skuIDs) == 0 {
		return []wmsmodel.InventoryStock{}, nil
	}
	var rows []wmsmodel.InventoryStock
	if err := db.DB.WithContext(ctx).
		Where("sku_id IN ?", skuIDs).
		Find(&rows).Error; err != nil {
		return nil, WrapDBError("查询库存失败", err)
	}
	return rows, nil
}

func UpdateStockSafety(ctx context.Context, id uint64, safeQty int) error {
	if id == 0 {
		return InvalidError("库存ID不能为空")
	}
	if safeQty < 0 {
		return InvalidError("安全库存不能小于0")
	}
	var row wmsmodel.InventoryStock
	if err := db.DB.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFoundError("库存记录不存在")
		}
		return WrapDBError("查询库存失败", err)
	}
	if err := db.DB.WithContext(ctx).Model(&wmsmodel.InventoryStock{}).Where("id = ?", id).Update("safe_qty", safeQty).Error; err != nil {
		return WrapDBError("更新安全库存失败", err)
	}
	return nil
}
