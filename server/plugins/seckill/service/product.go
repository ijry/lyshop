package service

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/server/core/db"
	seckillmodel "github.com/ijry/lyshop/server/plugins/seckill/model"
	"gorm.io/gorm"
)

// ListProducts 获取秒杀商品列表（管理后台）
func ListProducts(ctx context.Context, activityID uint64, page, size int) ([]seckillmodel.SeckillProduct, int64, error) {
	var products []seckillmodel.SeckillProduct
	var total int64

	query := db.DB.WithContext(ctx).Model(&seckillmodel.SeckillProduct{})

	if activityID > 0 {
		query = query.Where("activity_id = ?", activityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Order("sort DESC, id DESC").Offset(offset).Limit(size).Find(&products).Error

	return products, total, err
}

// ListFrontProducts 获取秒杀商品列表（前端）
func ListFrontProducts(ctx context.Context, activityID uint64, page, size int) ([]map[string]interface{}, int64, error) {
	var total int64
	query := db.DB.WithContext(ctx).Table("seckill_products sp").
		Select("sp.*, p.title, p.cover, p.price as original_price, ps.price as sku_price").
		Joins("LEFT JOIN products p ON sp.product_id = p.id").
		Joins("LEFT JOIN product_skus ps ON sp.sku_id = ps.id")

	if activityID > 0 {
		query = query.Where("sp.activity_id = ?", activityID)
	}

	// 只显示有效活动的商品
	query = query.Where("EXISTS (SELECT 1 FROM seckill_activities sa WHERE sa.id = sp.activity_id AND sa.status = 1 AND sa.start_at <= NOW() AND sa.end_at >= NOW())")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var results []map[string]interface{}
	offset := (page - 1) * size
	err := query.Order("sp.sort DESC, sp.id DESC").Offset(offset).Limit(size).Find(&results).Error

	return results, total, err
}

// GetProduct 获取秒杀商品详情
func GetProduct(ctx context.Context, id uint64) (*seckillmodel.SeckillProduct, error) {
	var product seckillmodel.SeckillProduct
	err := db.DB.WithContext(ctx).First(&product, id).Error
	return &product, err
}

// UpsertProducts 批量更新秒杀商品
func UpsertProducts(ctx context.Context, activityID uint64, products []seckillmodel.SeckillProduct) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除现有商品
		if err := tx.Where("activity_id = ?", activityID).Delete(&seckillmodel.SeckillProduct{}).Error; err != nil {
			return err
		}

		// 插入新商品
		if len(products) > 0 {
			for i := range products {
				products[i].ActivityID = activityID
			}
			if err := tx.Create(&products).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// IncreaseSoldQty 增加已售数量（事务中调用）
func IncreaseSoldQtyTx(tx *gorm.DB, productID uint64, qty int) error {
	return tx.Model(&seckillmodel.SeckillProduct{}).
		Where("id = ?", productID).
		Update("sold_qty", gorm.Expr("sold_qty + ?", qty)).Error
}

// ValidateProduct 验证秒杀商品
func ValidateProduct(ctx context.Context, productID uint64) (*seckillmodel.SeckillProduct, error) {
	var product seckillmodel.SeckillProduct
	if err := db.DB.WithContext(ctx).First(&product, productID).Error; err != nil {
		return nil, fmt.Errorf("秒杀商品不存在")
	}

	// 检查库存
	if product.TotalStockLimit > 0 && product.SoldQty >= product.TotalStockLimit {
		return nil, fmt.Errorf("秒杀商品已售罄")
	}

	// 检查活动是否有效
	var activity seckillmodel.SeckillActivity
	if err := db.DB.WithContext(ctx).First(&activity, product.ActivityID).Error; err != nil {
		return nil, fmt.Errorf("秒杀活动不存在")
	}

	if activity.Status != 1 {
		return nil, fmt.Errorf("秒杀活动已停用")
	}

	// 检查活动时间
	now := db.DB.NowFunc()
	if activity.StartAt != nil && now.Before(*activity.StartAt) {
		return nil, fmt.Errorf("秒杀活动未开始")
	}
	if activity.EndAt != nil && now.After(*activity.EndAt) {
		return nil, fmt.Errorf("秒杀活动已结束")
	}

	return &product, nil
}
