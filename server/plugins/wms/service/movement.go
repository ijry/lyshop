package service

import (
	"context"
	"strings"

	"github.com/ijry/lyshop/core/db"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
)

func ListMovements(ctx context.Context, q MovementListQuery) ([]wmsmodel.InventoryMovement, int64, error) {
	page, size := normalizePage(q.Page, q.Size)
	tx := db.DB.WithContext(ctx).Model(&wmsmodel.InventoryMovement{})
	if q.WarehouseID > 0 {
		tx = tx.Where("warehouse_id = ?", q.WarehouseID)
	}
	if q.SkuID > 0 {
		tx = tx.Where("sku_id = ?", q.SkuID)
	}
	if q.BizType != "" {
		tx = tx.Where("biz_type = ?", q.BizType)
	}
	if q.DocNo != "" {
		tx = tx.Where("doc_no = ?", strings.TrimSpace(q.DocNo))
	}
	if q.StartAt != nil {
		tx = tx.Where("occurred_at >= ?", *q.StartAt)
	}
	if q.EndAt != nil {
		tx = tx.Where("occurred_at <= ?", *q.EndAt)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, WrapDBError("查询流水总数失败", err)
	}
	var rows []wmsmodel.InventoryMovement
	if err := tx.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, 0, WrapDBError("查询流水列表失败", err)
	}
	return rows, total, nil
}
