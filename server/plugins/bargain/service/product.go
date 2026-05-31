package service

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/plugins/bargain/model"
)

// ListProducts 获取商品列表（管理端）
func ListProducts(ctx context.Context, activityID uint64, page, size int) ([]model.BargainProduct, int64, error) {
	var products []model.BargainProduct
	var total int64

	query := db.DB.WithContext(ctx).Model(&model.BargainProduct{})
	if activityID > 0 {
		query = query.Where("activity_id = ?", activityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("sort DESC, id DESC").Offset(offset).Limit(size).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListFrontProducts 获取商品列表（前端）
func ListFrontProducts(ctx context.Context, activityID uint64, page, size int) ([]model.BargainProduct, int64, error) {
	var products []model.BargainProduct
	var total int64

	query := db.DB.WithContext(ctx).
		Joins("JOIN bargain_activities ON bargain_activities.id = bargain_products.activity_id").
		Where("bargain_activities.status = 1").
		Where("bargain_activities.start_at <= NOW() AND bargain_activities.end_at >= NOW()")

	if activityID > 0 {
		query = query.Where("bargain_products.activity_id = ?", activityID)
	}

	if err := query.Model(&model.BargainProduct{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("bargain_products.sort DESC, bargain_products.id DESC").
		Offset(offset).Limit(size).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProduct 获取商品详情
func GetProduct(ctx context.Context, id uint64) (*model.BargainProduct, error) {
	var product model.BargainProduct
	if err := db.DB.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// UpsertProducts 批量更新商品
func UpsertProducts(ctx context.Context, activityID uint64, products []model.BargainProduct) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *db.Tx) error {
		// 删除旧商品
		if err := tx.Where("activity_id = ?", activityID).Delete(&model.BargainProduct{}).Error; err != nil {
			return err
		}

		// 插入新商品
		for i := range products {
			products[i].ActivityID = activityID
			if err := tx.Create(&products[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// IncreaseSoldQtyTx 增加已售数量（事务中）
func IncreaseSoldQtyTx(tx *db.Tx, productID uint64, qty int) error {
	return tx.Model(&model.BargainProduct{}).
		Where("id = ?", productID).
		UpdateColumn("sold_qty", db.Expr("sold_qty + ?", qty)).Error
}

// ValidateProduct 验证商品是否可购买
func ValidateProduct(ctx context.Context, productID uint64) error {
	product, err := GetProduct(ctx, productID)
	if err != nil {
		return err
	}

	// 检查库存
	if product.TotalStockLimit > 0 && product.SoldQty >= product.TotalStockLimit {
		return db.ErrStockInsufficient
	}

	// 检查活动有效性
	activity, err := GetActivity(ctx, product.ActivityID)
	if err != nil {
		return err
	}

	if activity.Status != 1 {
		return db.ErrActivityInactive
	}

	return nil
}
