package service

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/plugins/bargain/model"
)

var (
	ErrAlreadyHelped = errors.New("已帮助过该砍价")
	ErrMaxHelpers    = errors.New("助力人数已达上限")
	ErrBargainExpired = errors.New("砍价已过期")
	ErrCannotHelpSelf = errors.New("不能帮自己砍价")
)

// CreateBargainOrder 创建砍价订单（发起砍价）
func CreateBargainOrder(ctx context.Context, activityID, productID, skuID, userID uint64) (*model.BargainOrder, error) {
	product, err := GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	expireAt := time.Now().Add(time.Duration(product.ExpireHours) * time.Hour)

	bargainOrder := &model.BargainOrder{
		ActivityID:   activityID,
		ProductID:    productID,
		SkuID:        skuID,
		UserID:       userID,
		StartPrice:   product.StartPrice,
		FloorPrice:   product.FloorPrice,
		CurrentPrice: product.StartPrice,
		HelpCount:    0,
		Status:       "pending",
		ExpireAt:     &expireAt,
	}

	if err := db.DB.WithContext(ctx).Create(bargainOrder).Error; err != nil {
		return nil, err
	}

	return bargainOrder, nil
}

// HelpBargain 帮助砍价
func HelpBargain(ctx context.Context, bargainOrderID, helperUserID uint64) (float64, error) {
	var cutAmount float64

	err := db.DB.WithContext(ctx).Transaction(func(tx *db.Tx) error {
		var bargainOrder model.BargainOrder
		if err := tx.Where("id = ?", bargainOrderID).First(&bargainOrder).Error; err != nil {
			return err
		}

		// 检查状态
		if bargainOrder.Status != "pending" {
			return errors.New("砍价已结束")
		}

		// 检查是否过期
		if bargainOrder.ExpireAt != nil && time.Now().After(*bargainOrder.ExpireAt) {
			return ErrBargainExpired
		}

		// 不能帮自己砍价
		if bargainOrder.UserID == helperUserID {
			return ErrCannotHelpSelf
		}

		// 检查是否已帮助过
		var count int64
		if err := tx.Model(&model.BargainHelper{}).
			Where("bargain_order_id = ? AND user_id = ?", bargainOrderID, helperUserID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrAlreadyHelped
		}

		// 获取商品配置
		var product model.BargainProduct
		if err := tx.First(&product, bargainOrder.ProductID).Error; err != nil {
			return err
		}

		// 检查助力人数上限
		if product.MaxHelpers > 0 && bargainOrder.HelpCount >= product.MaxHelpers {
			return ErrMaxHelpers
		}

		// 计算砍价金额
		cutAmount = calculateCutAmount(
			bargainOrder.CurrentPrice,
			bargainOrder.FloorPrice,
			product.MinCutAmount,
			product.MaxCutAmount,
			bargainOrder.HelpCount,
			product.MaxHelpers,
		)

		// 确保不低于底价
		newPrice := bargainOrder.CurrentPrice - cutAmount
		if newPrice < bargainOrder.FloorPrice {
			cutAmount = bargainOrder.CurrentPrice - bargainOrder.FloorPrice
			newPrice = bargainOrder.FloorPrice
		}

		// 添加助力记录
		helper := &model.BargainHelper{
			BargainOrderID: bargainOrderID,
			UserID:         helperUserID,
			CutAmount:      cutAmount,
		}
		if err := tx.Create(helper).Error; err != nil {
			return err
		}

		// 更新砍价订单
		bargainOrder.CurrentPrice = newPrice
		bargainOrder.HelpCount++

		// 检查是否砍到底价
		if newPrice <= bargainOrder.FloorPrice {
			bargainOrder.Status = "success"
		}

		if err := tx.Save(&bargainOrder).Error; err != nil {
			return err
		}

		return nil
	})

	return cutAmount, err
}

// calculateCutAmount 计算砍价金额（使用随机算法）
func calculateCutAmount(currentPrice, floorPrice, minCut, maxCut float64, helpCount, maxHelpers int) float64 {
	remaining := currentPrice - floorPrice
	if remaining <= 0 {
		return 0
	}

	// 如果是最后一次助力，直接砍到底价
	if maxHelpers > 0 && helpCount >= maxHelpers-1 {
		return remaining
	}

	// 随机砍价金额，在最小和最大之间
	cutAmount := minCut + rand.Float64()*(maxCut-minCut)

	// 确保不超过剩余金额
	if cutAmount > remaining {
		cutAmount = remaining
	}

	// 保留两位小数
	return math.Round(cutAmount*100) / 100
}

// ExpireBargain 砍价过期处理
func ExpireBargain(ctx context.Context, bargainOrderID uint64) error {
	return db.DB.WithContext(ctx).Model(&model.BargainOrder{}).
		Where("id = ?", bargainOrderID).
		Where("status = ?", "pending").
		Update("status", "failed").Error
}

// ListBargainOrders 获取砍价订单列表
func ListBargainOrders(ctx context.Context, page, size int, status string) ([]model.BargainOrder, int64, error) {
	var orders []model.BargainOrder
	var total int64

	query := db.DB.WithContext(ctx).Model(&model.BargainOrder{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("id DESC").Offset(offset).Limit(size).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetBargainOrder 获取砍价订单详情
func GetBargainOrder(ctx context.Context, id uint64) (*model.BargainOrder, error) {
	var order model.BargainOrder
	if err := db.DB.WithContext(ctx).First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// ListBargainHelpers 获取助力记录列表
func ListBargainHelpers(ctx context.Context, bargainOrderID uint64) ([]model.BargainHelper, error) {
	var helpers []model.BargainHelper
	err := db.DB.WithContext(ctx).
		Where("bargain_order_id = ?", bargainOrderID).
		Order("id ASC").
		Find(&helpers).Error
	return helpers, err
}

// CheckExpiredBargains 检查并处理过期砍价（定时任务）
func CheckExpiredBargains(ctx context.Context) error {
	now := time.Now()
	return db.DB.WithContext(ctx).Model(&model.BargainOrder{}).
		Where("status = ?", "pending").
		Where("expire_at < ?", now).
		Update("status", "failed").Error
}
