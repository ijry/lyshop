package inventory

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type localProvider struct{}

func init() {
	Register(&localProvider{})
}

func (p *localProvider) Name() string { return "local" }

func (p *localProvider) ReserveTx(tx *gorm.DB, in ReserveInput) error {
	for _, item := range in.Items {
		var sku productmodel.ProductSku
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&sku, item.SkuID).Error; err != nil {
			return err
		}
		if sku.Status != productmodel.ProductSkuStatusActive {
			return fmt.Errorf("SKU已下线")
		}
		if sku.Stock < item.Qty {
			return fmt.Errorf("库存不足")
		}
		if err := tx.Create(&InventoryReservation{
			BizType:   in.BizType,
			BizNo:     in.BizNo,
			SkuID:     item.SkuID,
			Qty:       item.Qty,
			Status:    InventoryStatusReserved,
			ExpiredAt: in.ExpiredAt,
		}).Error; err != nil {
			return err
		}
	}
	return tx.Where("order_no = ?", in.BizNo).
		Assign(&OrderInventoryState{
			BizType:  in.BizType,
			OrderNo:  in.BizNo,
			Status:   InventoryStatusReserved,
			Provider: p.Name(),
		}).
		FirstOrCreate(&OrderInventoryState{}).Error
}

func (p *localProvider) ConfirmTx(tx *gorm.DB, bizType, bizNo string) error {
	var rows []InventoryReservation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("biz_type = ? AND biz_no = ?", bizType, bizNo).
		Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if row.Status == InventoryStatusConfirmed {
			continue
		}
		if row.Status != InventoryStatusReserved {
			return fmt.Errorf("inventory reservation status invalid: %s", row.Status)
		}
		result := tx.Model(&productmodel.ProductSku{}).
			Where("id = ? AND stock >= ?", row.SkuID, row.Qty).
			Update("stock", gorm.Expr("stock - ?", row.Qty))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("库存不足")
		}
		if err := tx.Model(&InventoryReservation{}).
			Where("id = ?", row.ID).
			Update("status", InventoryStatusConfirmed).Error; err != nil {
			return err
		}
	}
	return tx.Model(&OrderInventoryState{}).
		Where("order_no = ?", bizNo).
		Updates(map[string]any{
			"status":     InventoryStatusConfirmed,
			"last_error": "",
		}).Error
}

func (p *localProvider) ReleaseTx(tx *gorm.DB, bizType, bizNo, reason string) error {
	var rows []InventoryReservation
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("biz_type = ? AND biz_no = ?", bizType, bizNo).
		Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if row.Status == InventoryStatusReleased {
			continue
		}
		if row.Status == InventoryStatusConfirmed {
			return fmt.Errorf("inventory reservation already confirmed")
		}
		if err := tx.Model(&InventoryReservation{}).
			Where("id = ?", row.ID).
			Updates(map[string]any{
				"status": InventoryStatusReleased,
				"reason": reason,
			}).Error; err != nil {
			return err
		}
	}
	return tx.Model(&OrderInventoryState{}).
		Where("order_no = ?", bizNo).
		Updates(map[string]any{
			"status":     InventoryStatusReleased,
			"last_error": reason,
		}).Error
}

func (p *localProvider) DeductTx(tx *gorm.DB, in DeductInput) error {
	for _, item := range in.Items {
		result := tx.Model(&productmodel.ProductSku{}).
			Where("id = ? AND stock >= ?", item.SkuID, item.Qty).
			Update("stock", gorm.Expr("stock - ?", item.Qty))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("库存不足")
		}
	}
	return nil
}

func (p *localProvider) RestoreTx(tx *gorm.DB, in RestoreInput) error {
	for _, item := range in.Items {
		if err := tx.Model(&productmodel.ProductSku{}).
			Where("id = ?", item.SkuID).
			Update("stock", gorm.Expr("stock + ?", item.Qty)).Error; err != nil {
			return err
		}
	}
	return nil
}

func (p *localProvider) SyncSkuTx(_ *gorm.DB, _ SyncSkuInput) error { return nil }

func (p *localProvider) GetSellableStock(ctx context.Context, skuIDs []uint64) ([]SellableStock, error) {
	var rows []productmodel.ProductSku
	if err := db.DB.WithContext(ctx).Where("id IN ?", skuIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]SellableStock, 0, len(rows))
	for _, row := range rows {
		out = append(out, SellableStock{
			SkuID:         row.ID,
			SellableStock: row.Stock,
			Reserved:      0,
			OnHand:        row.Stock,
		})
	}
	return out, nil
}
