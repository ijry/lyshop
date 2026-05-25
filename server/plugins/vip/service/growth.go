package service

import (
	"context"
	"fmt"
	"math"

	"github.com/ijry/lyshop/core/db"
	vipmodel "github.com/ijry/lyshop/plugins/vip/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func calcGrowth(amount float64) int64 {
	if amount <= 0 {
		return 0
	}
	return int64(math.Floor(amount))
}

func applyGrowthDeltaTx(ctx context.Context, tx *gorm.DB, userID, orderID uint64, eventType string, delta int64, idempotencyKey, remark string) error {
	if idempotencyKey == "" || delta == 0 {
		return nil
	}
	var existing vipmodel.GrowthLog
	if err := tx.WithContext(ctx).Where("idempotency_key = ?", idempotencyKey).First(&existing).Error; err == nil {
		return nil
	}

	var asset vipmodel.UserAsset
	if err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&asset).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		asset = vipmodel.UserAsset{UserID: userID, Status: 0, GrowthValue: 0}
		if err := tx.WithContext(ctx).Create(&asset).Error; err != nil {
			return err
		}
	}
	next := asset.GrowthValue + delta
	if next < 0 {
		next = 0
	}

	levelID := asset.CurrentLevelID
	var levels []vipmodel.Level
	tx.WithContext(ctx).Where("status = 1").Order("growth_threshold asc, id asc").Find(&levels)
	for _, lv := range levels {
		if next >= lv.GrowthThreshold {
			levelID = lv.ID
		}
	}

	if err := tx.WithContext(ctx).Model(&vipmodel.UserAsset{}).Where("id = ?", asset.ID).Updates(map[string]any{
		"growth_value":     next,
		"current_level_id": levelID,
	}).Error; err != nil {
		return err
	}

	return tx.WithContext(ctx).Create(&vipmodel.GrowthLog{
		UserID:         userID,
		OrderID:        orderID,
		EventType:      eventType,
		GrowthDelta:    delta,
		BalanceAfter:   next,
		IdempotencyKey: idempotencyKey,
		Remark:         remark,
	}).Error
}

func GrantGrowthForPaidOrderTx(tx *gorm.DB, userID, orderID uint64, payable float64) error {
	delta := calcGrowth(payable)
	if delta <= 0 {
		return nil
	}
	key := fmt.Sprintf("vip:growth:paid:%d", orderID)
	ctx := tx.Statement.Context
	if err := applyGrowthDeltaTx(ctx, tx, userID, orderID, "order_paid", delta, key, "订单支付成长值"); err != nil {
		return err
	}
	return tx.WithContext(ctx).Model(&vipmodel.OrderBenefit{}).
		Where(vipmodel.OrderBenefit{OrderID: orderID}).
		Assign(vipmodel.OrderBenefit{UserID: userID, GrowthGranted: delta}).
		FirstOrCreate(&vipmodel.OrderBenefit{}).Error
}

func RollbackGrowthForRefundTx(tx *gorm.DB, userID, orderID, caseID uint64, refundAmount float64) error {
	delta := calcGrowth(refundAmount)
	if delta <= 0 {
		return nil
	}
	key := fmt.Sprintf("vip:growth:refund:case:%d", caseID)
	ctx := tx.Statement.Context
	if err := applyGrowthDeltaTx(ctx, tx, userID, orderID, "order_refund", -delta, key, "订单退款回退成长值"); err != nil {
		return err
	}
	return tx.WithContext(ctx).Model(&vipmodel.OrderBenefit{}).
		Where(vipmodel.OrderBenefit{OrderID: orderID}).
		Assign(map[string]any{"user_id": userID, "growth_reverted": gorm.Expr("growth_reverted + ?", delta)}).
		FirstOrCreate(&vipmodel.OrderBenefit{}).Error
}

func GrantGrowthOnOrderPaid(ctx context.Context, orderID, userID uint64, payable float64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return GrantGrowthForPaidOrderTx(tx, userID, orderID, payable)
	})
}

func RevertGrowthOnRefund(ctx context.Context, orderID, userID, caseID uint64, refundAmount float64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return RollbackGrowthForRefundTx(tx, userID, orderID, caseID, refundAmount)
	})
}
