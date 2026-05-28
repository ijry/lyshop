package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DocDetail struct {
	Doc   wmsmodel.InventoryDoc       `json:"doc"`
	Items []wmsmodel.InventoryDocItem `json:"items"`
}

func ListDocs(ctx context.Context, q DocListQuery) ([]wmsmodel.InventoryDoc, int64, error) {
	page, size := normalizePage(q.Page, q.Size)
	tx := db.DB.WithContext(ctx).Model(&wmsmodel.InventoryDoc{})
	if q.WarehouseID > 0 {
		tx = tx.Where("warehouse_id = ?", q.WarehouseID)
	}
	if q.DocType != "" {
		tx = tx.Where("doc_type = ?", q.DocType)
	}
	if q.Status != "" {
		tx = tx.Where("status = ?", q.Status)
	}
	if q.DocNo != "" {
		tx = tx.Where("doc_no LIKE ?", "%"+strings.TrimSpace(q.DocNo)+"%")
	}
	if q.StartAt != nil {
		tx = tx.Where("created_at >= ?", *q.StartAt)
	}
	if q.EndAt != nil {
		tx = tx.Where("created_at <= ?", *q.EndAt)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []wmsmodel.InventoryDoc
	if err := tx.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func CreateDraftDoc(ctx context.Context, in CreateDocInput) (*wmsmodel.InventoryDoc, error) {
	if !wmsmodel.IsValidDocType(in.DocType) {
		return nil, fmt.Errorf("单据类型非法")
	}
	if len(in.Items) == 0 {
		return nil, fmt.Errorf("单据明细不能为空")
	}
	if in.WarehouseID == 0 {
		return nil, fmt.Errorf("仓库ID不能为空")
	}
	if err := validateDocItems(in.Items); err != nil {
		return nil, err
	}

	doc := &wmsmodel.InventoryDoc{
		DocNo:       genDocNo(in.DocType),
		DocType:     in.DocType,
		Status:      wmsmodel.DocStatusDraft,
		WarehouseID: in.WarehouseID,
		Remark:      strings.TrimSpace(in.Remark),
	}
	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensureWarehouseExists(tx, in.WarehouseID); err != nil {
			return err
		}
		if err := tx.Create(doc).Error; err != nil {
			return err
		}
		items := make([]wmsmodel.InventoryDocItem, 0, len(in.Items))
		for _, item := range in.Items {
			items = append(items, wmsmodel.InventoryDocItem{
				DocID:  doc.ID,
				SkuID:  item.SkuID,
				Qty:    item.Qty,
				Remark: strings.TrimSpace(item.Remark),
			})
		}
		return tx.Create(&items).Error
	})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func GetDocDetail(ctx context.Context, id uint64) (*DocDetail, error) {
	if id == 0 {
		return nil, fmt.Errorf("单据ID不能为空")
	}
	var doc wmsmodel.InventoryDoc
	if err := db.DB.WithContext(ctx).Where("id = ?", id).First(&doc).Error; err != nil {
		return nil, fmt.Errorf("单据不存在")
	}
	var items []wmsmodel.InventoryDocItem
	if err := db.DB.WithContext(ctx).Where("doc_id = ?", id).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return &DocDetail{Doc: doc, Items: items}, nil
}

func UpdateDraftDoc(ctx context.Context, id uint64, in UpdateDocInput) error {
	if id == 0 {
		return fmt.Errorf("单据ID不能为空")
	}
	if !wmsmodel.IsValidDocType(in.DocType) {
		return fmt.Errorf("单据类型非法")
	}
	if len(in.Items) == 0 {
		return fmt.Errorf("单据明细不能为空")
	}
	if in.WarehouseID == 0 {
		return fmt.Errorf("仓库ID不能为空")
	}
	if err := validateDocItems(in.Items); err != nil {
		return err
	}

	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		doc, err := lockDoc(tx, id)
		if err != nil {
			return err
		}
		if doc.Status != wmsmodel.DocStatusDraft {
			return fmt.Errorf("单据状态非法，只有草稿单允许编辑")
		}
		if err := ensureWarehouseExists(tx, in.WarehouseID); err != nil {
			return err
		}
		if err := tx.Model(&wmsmodel.InventoryDoc{}).Where("id = ?", id).Updates(map[string]any{
			"warehouse_id": in.WarehouseID,
			"doc_type":     in.DocType,
			"remark":       strings.TrimSpace(in.Remark),
		}).Error; err != nil {
			return err
		}
		if err := tx.Where("doc_id = ?", id).Delete(&wmsmodel.InventoryDocItem{}).Error; err != nil {
			return err
		}
		items := make([]wmsmodel.InventoryDocItem, 0, len(in.Items))
		for _, item := range in.Items {
			items = append(items, wmsmodel.InventoryDocItem{
				DocID:  id,
				SkuID:  item.SkuID,
				Qty:    item.Qty,
				Remark: strings.TrimSpace(item.Remark),
			})
		}
		return tx.Create(&items).Error
	})
}

func CancelDraftDoc(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("单据ID不能为空")
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		doc, err := lockDoc(tx, id)
		if err != nil {
			return err
		}
		if doc.Status != wmsmodel.DocStatusDraft {
			return fmt.Errorf("单据状态非法，只有草稿单允许作废")
		}
		now := time.Now()
		return tx.Model(&wmsmodel.InventoryDoc{}).Where("id = ?", id).Updates(map[string]any{
			"status":      wmsmodel.DocStatusCanceled,
			"canceled_at": now,
		}).Error
	})
}

func CompleteDraftDoc(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("单据ID不能为空")
	}

	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		doc, err := lockDoc(tx, id)
		if err != nil {
			return err
		}
		if doc.Status != wmsmodel.DocStatusDraft {
			return fmt.Errorf("单据状态非法，只有草稿单允许完成")
		}

		warehouse, err := lockWarehouse(tx, doc.WarehouseID)
		if err != nil {
			return err
		}
		if warehouse.Status != wmsmodel.WarehouseStatusEnabled {
			return fmt.Errorf("仓库停用，不能完成单据")
		}

		var items []wmsmodel.InventoryDocItem
		if err := tx.Where("doc_id = ?", doc.ID).Order("id ASC").Find(&items).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return fmt.Errorf("单据明细不能为空")
		}

		for _, item := range items {
			stock, err := lockOrCreateStock(tx, doc, item)
			if err != nil {
				return err
			}
			before := stock.Qty
			after := before
			change := item.Qty
			if doc.DocType == wmsmodel.DocTypeOutbound {
				if before < item.Qty {
					return fmt.Errorf("库存不足，SKU=%d 当前=%d 需要=%d", item.SkuID, before, item.Qty)
				}
				change = -item.Qty
			}
			after = before + change
			if err := tx.Model(&wmsmodel.InventoryStock{}).Where("id = ?", stock.ID).Update("qty", after).Error; err != nil {
				return err
			}
			mv := wmsmodel.InventoryMovement{
				WarehouseID: doc.WarehouseID,
				SkuID:       item.SkuID,
				DocID:       doc.ID,
				DocNo:       doc.DocNo,
				BizType:     doc.DocType,
				ChangeQty:   change,
				BeforeQty:   before,
				AfterQty:    after,
				OccurredAt:  time.Now(),
				Remark:      strings.TrimSpace(item.Remark),
			}
			if err := tx.Create(&mv).Error; err != nil {
				return err
			}
		}
		now := time.Now()
		return tx.Model(&wmsmodel.InventoryDoc{}).Where("id = ?", doc.ID).Updates(map[string]any{
			"status":       wmsmodel.DocStatusCompleted,
			"completed_at": now,
		}).Error
	})
}

func lockDoc(tx *gorm.DB, id uint64) (*wmsmodel.InventoryDoc, error) {
	var doc wmsmodel.InventoryDoc
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&doc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("单据不存在")
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func lockWarehouse(tx *gorm.DB, id uint64) (*wmsmodel.Warehouse, error) {
	var row wmsmodel.Warehouse
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("仓库不存在")
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func ensureWarehouseExists(tx *gorm.DB, warehouseID uint64) error {
	var cnt int64
	if err := tx.Model(&wmsmodel.Warehouse{}).Where("id = ?", warehouseID).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		return fmt.Errorf("仓库不存在")
	}
	return nil
}

func validateDocItems(items []DocItemInput) error {
	for _, item := range items {
		if item.SkuID == 0 {
			return fmt.Errorf("SKU不能为空")
		}
		if item.Qty <= 0 {
			return fmt.Errorf("数量必须大于0")
		}
	}
	return nil
}

func lockOrCreateStock(tx *gorm.DB, doc *wmsmodel.InventoryDoc, item wmsmodel.InventoryDocItem) (*wmsmodel.InventoryStock, error) {
	var stock wmsmodel.InventoryStock
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("warehouse_id = ? AND sku_id = ?", doc.WarehouseID, item.SkuID).
		First(&stock).Error
	if err == nil {
		return &stock, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if doc.DocType == wmsmodel.DocTypeOutbound {
		return nil, fmt.Errorf("库存不足，SKU=%d 当前=0 需要=%d", item.SkuID, item.Qty)
	}
	newStock := wmsmodel.InventoryStock{WarehouseID: doc.WarehouseID, SkuID: item.SkuID}
	if err := tx.Where("warehouse_id = ? AND sku_id = ?", doc.WarehouseID, item.SkuID).FirstOrCreate(&newStock).Error; err != nil {
		return nil, err
	}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("warehouse_id = ? AND sku_id = ?", doc.WarehouseID, item.SkuID).
		First(&stock).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}
