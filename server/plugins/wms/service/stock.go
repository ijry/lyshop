package service

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/core/db"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Inbound adds qty to stock for a given warehouse+SKU and records a log.
func Inbound(ctx context.Context, warehouseID, skuID uint64, qty int, refID uint64, refType string) error {
	if qty <= 0 {
		return fmt.Errorf("inbound qty must be positive")
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Upsert stock record
		stock := wmsmodel.WmsStock{WarehouseID: warehouseID, SkuID: skuID}
		tx.Where("warehouse_id = ? AND sku_id = ?", warehouseID, skuID).FirstOrCreate(&stock)

		before := stock.Qty
		after := before + qty
		if err := tx.Model(&stock).UpdateColumn("qty", after).Error; err != nil {
			return err
		}
		return tx.Create(&wmsmodel.WmsStockLog{
			WarehouseID: warehouseID, SkuID: skuID,
			Type: "inbound", Qty: qty, BeforeQty: before, AfterQty: after,
			RefID: refID, RefType: refType,
		}).Error
	})
}

// Outbound deducts qty from stock. Returns error if insufficient.
func Outbound(ctx context.Context, warehouseID, skuID uint64, qty int, refID uint64, refType string) error {
	if qty <= 0 {
		return fmt.Errorf("outbound qty must be positive")
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stock wmsmodel.WmsStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND sku_id = ?", warehouseID, skuID).
			First(&stock).Error; err != nil {
			return fmt.Errorf("库存记录不存在")
		}
		if stock.Qty < qty {
			return fmt.Errorf("库存不足 (当前=%d, 需要=%d)", stock.Qty, qty)
		}
		before := stock.Qty
		after := before - qty
		tx.Model(&stock).UpdateColumn("qty", after)
		return tx.Create(&wmsmodel.WmsStockLog{
			WarehouseID: warehouseID, SkuID: skuID,
			Type: "outbound", Qty: -qty, BeforeQty: before, AfterQty: after,
			RefID: refID, RefType: refType,
		}).Error
	})
}

// ListStocks returns paginated stock records for a warehouse.
func ListStocks(ctx context.Context, warehouseID uint64, page, size int) ([]wmsmodel.WmsStock, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 100 { size = 20 }
	tx := db.DB.WithContext(ctx).Model(&wmsmodel.WmsStock{})
	if warehouseID > 0 {
		tx = tx.Where("warehouse_id = ?", warehouseID)
	}
	var total int64
	tx.Count(&total)
	var list []wmsmodel.WmsStock
	err := tx.Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

// ListWarehouses returns all active warehouses.
func ListWarehouses(ctx context.Context) ([]wmsmodel.Warehouse, error) {
	var list []wmsmodel.Warehouse
	err := db.DB.WithContext(ctx).Where("status = 1").Find(&list).Error
	return list, err
}

func CreateWarehouse(ctx context.Context, w *wmsmodel.Warehouse) error {
	return db.DB.WithContext(ctx).Create(w).Error
}
