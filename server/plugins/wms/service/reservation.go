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

const (
	MovementBizTypeReserve       = "reserve"
	MovementBizTypeRelease       = "release"
	MovementBizTypeOrderOutbound = "order_outbound"
)

type ReservationItemInput struct {
	SkuID  uint64
	Qty    int
	Remark string
}

type ReserveStockInput struct {
	BizType     string
	BizNo       string
	WarehouseID uint64
	Items       []ReservationItemInput
	ExpiredAt   *time.Time
	Remark      string
}

func PickDefaultWarehouseID(ctx context.Context) (uint64, error) {
	var warehouse wmsmodel.Warehouse
	if err := db.DB.WithContext(ctx).
		Where("status = ?", wmsmodel.WarehouseStatusEnabled).
		Order("id ASC").
		First(&warehouse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, NotFoundError("无可用仓库")
		}
		return 0, WrapDBError("查询默认仓库失败", err)
	}
	return warehouse.ID, nil
}

func PickDefaultWarehouseIDTx(tx *gorm.DB) (uint64, error) {
	var warehouse wmsmodel.Warehouse
	if err := tx.
		Where("status = ?", wmsmodel.WarehouseStatusEnabled).
		Order("id ASC").
		First(&warehouse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, NotFoundError("无可用仓库")
		}
		return 0, WrapDBError("查询默认仓库失败", err)
	}
	return warehouse.ID, nil
}

func ReserveStock(ctx context.Context, in ReserveStockInput) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return ReserveStockTx(tx, in)
	})
}

func ReserveStockTx(tx *gorm.DB, in ReserveStockInput) error {
	bizType := strings.TrimSpace(in.BizType)
	bizNo := strings.TrimSpace(in.BizNo)
	if bizType == "" || bizNo == "" {
		return InvalidError("业务类型和业务单号不能为空")
	}
	if in.WarehouseID == 0 {
		return InvalidError("仓库ID不能为空")
	}
	items, err := normalizeReservationItems(in.Items)
	if err != nil {
		return err
	}

	var existing []wmsmodel.InventoryReservation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("biz_type = ? AND biz_no = ?", bizType, bizNo).
		Find(&existing).Error; err != nil {
		return WrapDBError("查询预占记录失败", err)
	}
	if len(existing) > 0 {
		return ensureReserveIdempotent(existing, in.WarehouseID, items)
	}

	now := time.Now()
	for _, item := range items {
		var stock wmsmodel.InventoryStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND sku_id = ?", in.WarehouseID, item.SkuID).
			First(&stock).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ConflictError(fmt.Sprintf("库存不足，SKU=%d 当前=0 需要=%d", item.SkuID, item.Qty))
			}
			return WrapDBError("锁定库存失败", err)
		}
		sellable := stock.Qty - stock.ReservedQty
		if sellable < item.Qty {
			return ConflictError(fmt.Sprintf("库存不足，SKU=%d 当前可售=%d 需要=%d", item.SkuID, sellable, item.Qty))
		}
		if err := tx.Model(&wmsmodel.InventoryStock{}).Where("id = ?", stock.ID).Update("reserved_qty", stock.ReservedQty+item.Qty).Error; err != nil {
			return WrapDBError("更新预占库存失败", err)
		}

		reservation := wmsmodel.InventoryReservation{
			BizType:     bizType,
			BizNo:       bizNo,
			WarehouseID: in.WarehouseID,
			SkuID:       item.SkuID,
			Qty:         item.Qty,
			Status:      wmsmodel.ReservationStatusReserved,
			ExpiredAt:   in.ExpiredAt,
			Remark:      strings.TrimSpace(item.Remark),
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return WrapDBError("创建预占记录失败", err)
		}

		mv := wmsmodel.InventoryMovement{
			WarehouseID: in.WarehouseID,
			SkuID:       item.SkuID,
			DocID:       0,
			DocNo:       bizNo,
			BizType:     MovementBizTypeReserve,
			ChangeQty:   0,
			BeforeQty:   stock.Qty,
			AfterQty:    stock.Qty,
			OccurredAt:  now,
			Remark:      fmt.Sprintf("biz_type=%s,reserved=%d", bizType, item.Qty),
		}
		if err := tx.Create(&mv).Error; err != nil {
			return WrapDBError("写入预占流水失败", err)
		}
	}
	return nil
}

func ConfirmReservation(ctx context.Context, bizType, bizNo string) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return ConfirmReservationTx(tx, bizType, bizNo)
	})
}

func ConfirmReservationTx(tx *gorm.DB, bizType, bizNo string) error {
	bizType = strings.TrimSpace(bizType)
	bizNo = strings.TrimSpace(bizNo)
	if bizType == "" || bizNo == "" {
		return InvalidError("业务类型和业务单号不能为空")
	}
	var rows []wmsmodel.InventoryReservation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("biz_type = ? AND biz_no = ?", bizType, bizNo).
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return WrapDBError("查询预占记录失败", err)
	}
	if len(rows) == 0 {
		return NotFoundError("预占记录不存在")
	}

	allConfirmed := true
	for _, row := range rows {
		if row.Status != wmsmodel.ReservationStatusConfirmed {
			allConfirmed = false
		}
	}
	if allConfirmed {
		return nil
	}

	now := time.Now()
	for _, row := range rows {
		if row.Status == wmsmodel.ReservationStatusConfirmed {
			continue
		}
		if row.Status != wmsmodel.ReservationStatusReserved {
			return ConflictError("预占状态非法，无法确认")
		}
		var stock wmsmodel.InventoryStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND sku_id = ?", row.WarehouseID, row.SkuID).
			First(&stock).Error; err != nil {
			return WrapDBError("锁定库存失败", err)
		}
		if stock.ReservedQty < row.Qty || stock.Qty < row.Qty {
			return ConflictError(fmt.Sprintf("确认预占失败，SKU=%d 预占或在手库存不足", row.SkuID))
		}
		if err := tx.Model(&wmsmodel.InventoryStock{}).Where("id = ?", stock.ID).Updates(map[string]any{
			"qty":          stock.Qty - row.Qty,
			"reserved_qty": stock.ReservedQty - row.Qty,
		}).Error; err != nil {
			return WrapDBError("确认预占更新库存失败", err)
		}
		if err := tx.Model(&wmsmodel.InventoryReservation{}).Where("id = ?", row.ID).Update("status", wmsmodel.ReservationStatusConfirmed).Error; err != nil {
			return WrapDBError("更新预占状态失败", err)
		}
		mv := wmsmodel.InventoryMovement{
			WarehouseID: row.WarehouseID,
			SkuID:       row.SkuID,
			DocID:       0,
			DocNo:       row.BizNo,
			BizType:     MovementBizTypeOrderOutbound,
			ChangeQty:   -row.Qty,
			BeforeQty:   stock.Qty,
			AfterQty:    stock.Qty - row.Qty,
			OccurredAt:  now,
			Remark:      fmt.Sprintf("biz_type=%s,confirm=%d", row.BizType, row.Qty),
		}
		if err := tx.Create(&mv).Error; err != nil {
			return WrapDBError("写入确认流水失败", err)
		}
	}
	return nil
}

func ReleaseReservation(ctx context.Context, bizType, bizNo, reason string) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return ReleaseReservationTx(tx, bizType, bizNo, reason)
	})
}

func ReleaseReservationTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	bizType = strings.TrimSpace(bizType)
	bizNo = strings.TrimSpace(bizNo)
	reason = strings.TrimSpace(reason)
	if bizType == "" || bizNo == "" {
		return InvalidError("业务类型和业务单号不能为空")
	}
	var rows []wmsmodel.InventoryReservation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("biz_type = ? AND biz_no = ?", bizType, bizNo).
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return WrapDBError("查询预占记录失败", err)
	}
	if len(rows) == 0 {
		return nil
	}
	allReleased := true
	for _, row := range rows {
		if row.Status != wmsmodel.ReservationStatusReleased {
			allReleased = false
		}
		if row.Status == wmsmodel.ReservationStatusConfirmed {
			return ConflictError("预占已确认出库，不能释放")
		}
	}
	if allReleased {
		return nil
	}

	now := time.Now()
	for _, row := range rows {
		if row.Status == wmsmodel.ReservationStatusReleased {
			continue
		}
		if row.Status != wmsmodel.ReservationStatusReserved {
			return ConflictError("预占状态非法，无法释放")
		}
		var stock wmsmodel.InventoryStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND sku_id = ?", row.WarehouseID, row.SkuID).
			First(&stock).Error; err != nil {
			return WrapDBError("锁定库存失败", err)
		}
		if stock.ReservedQty < row.Qty {
			return ConflictError(fmt.Sprintf("释放预占失败，SKU=%d 预占库存不足", row.SkuID))
		}
		if err := tx.Model(&wmsmodel.InventoryStock{}).Where("id = ?", stock.ID).Update("reserved_qty", stock.ReservedQty-row.Qty).Error; err != nil {
			return WrapDBError("释放预占更新库存失败", err)
		}
		if err := tx.Model(&wmsmodel.InventoryReservation{}).Where("id = ?", row.ID).Updates(map[string]any{
			"status": wmsmodel.ReservationStatusReleased,
			"remark": reason,
		}).Error; err != nil {
			return WrapDBError("更新预占状态失败", err)
		}
		mv := wmsmodel.InventoryMovement{
			WarehouseID: row.WarehouseID,
			SkuID:       row.SkuID,
			DocID:       0,
			DocNo:       row.BizNo,
			BizType:     MovementBizTypeRelease,
			ChangeQty:   0,
			BeforeQty:   stock.Qty,
			AfterQty:    stock.Qty,
			OccurredAt:  now,
			Remark:      fmt.Sprintf("biz_type=%s,release=%d,reason=%s", row.BizType, row.Qty, reason),
		}
		if err := tx.Create(&mv).Error; err != nil {
			return WrapDBError("写入释放流水失败", err)
		}
	}
	return nil
}

func normalizeReservationItems(items []ReservationItemInput) ([]ReservationItemInput, error) {
	if len(items) == 0 {
		return nil, InvalidError("预占明细不能为空")
	}
	m := make(map[uint64]ReservationItemInput, len(items))
	for _, row := range items {
		if row.SkuID == 0 {
			return nil, InvalidError("SKU不能为空")
		}
		if row.Qty <= 0 {
			return nil, InvalidError("数量必须大于0")
		}
		existing, ok := m[row.SkuID]
		if ok {
			existing.Qty += row.Qty
			m[row.SkuID] = existing
			continue
		}
		m[row.SkuID] = ReservationItemInput{
			SkuID:  row.SkuID,
			Qty:    row.Qty,
			Remark: strings.TrimSpace(row.Remark),
		}
	}
	out := make([]ReservationItemInput, 0, len(m))
	for _, item := range m {
		out = append(out, item)
	}
	return out, nil
}

func ensureReserveIdempotent(existing []wmsmodel.InventoryReservation, warehouseID uint64, items []ReservationItemInput) error {
	expected := make(map[uint64]int, len(items))
	for _, item := range items {
		expected[item.SkuID] = item.Qty
	}
	if len(expected) != len(existing) {
		return ConflictError("预占记录已存在且与当前请求不一致")
	}
	for _, row := range existing {
		if row.WarehouseID != warehouseID {
			return ConflictError("预占记录已存在且仓库不一致")
		}
		qty, ok := expected[row.SkuID]
		if !ok || qty != row.Qty {
			return ConflictError("预占记录已存在且明细不一致")
		}
		if row.Status != wmsmodel.ReservationStatusReserved {
			return ConflictError("预占记录状态不允许重复预占")
		}
	}
	return nil
}
