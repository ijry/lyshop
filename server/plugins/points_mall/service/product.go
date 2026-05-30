package service

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/server/core/db"
	pmmodel "github.com/ijry/lyshop/server/plugins/points_mall/model"
)

// ListProducts 获取积分商品列表
func ListProducts(ctx context.Context, productType string, status int8, page, size int) ([]pmmodel.PointsProduct, int64, error) {
	var products []pmmodel.PointsProduct
	var total int64

	query := db.DB.WithContext(ctx).Model(&pmmodel.PointsProduct{})

	if productType != "" {
		query = query.Where("type = ?", productType)
	}

	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Order("sort DESC, id DESC").Offset(offset).Limit(size).Find(&products).Error

	return products, total, err
}

// GetProduct 获取积分商品详情
func GetProduct(ctx context.Context, id uint64) (*pmmodel.PointsProduct, error) {
	var product pmmodel.PointsProduct
	err := db.DB.WithContext(ctx).First(&product, id).Error
	return &product, err
}

// CreateProduct 创建积分商品
func CreateProduct(ctx context.Context, product *pmmodel.PointsProduct) error {
	return db.DB.WithContext(ctx).Create(product).Error
}

// UpdateProduct 更新积分商品
func UpdateProduct(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&pmmodel.PointsProduct{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// DeleteProduct 删除积分商品
func DeleteProduct(ctx context.Context, id uint64) error {
	// 检查是否有未完成的兑换记录
	var count int64
	db.DB.WithContext(ctx).Model(&pmmodel.PointsExchange{}).
		Where("product_id = ? AND status IN (?)", id, []string{"pending_ship", "shipped"}).
		Count(&count)

	if count > 0 {
		return fmt.Errorf("该商品有未完成的兑换记录，无法删除")
	}

	return db.DB.WithContext(ctx).Delete(&pmmodel.PointsProduct{}, id).Error
}

// UpdateProductStatus 更新商品状态（上下架）
func UpdateProductStatus(ctx context.Context, id uint64, status int8) error {
	return db.DB.WithContext(ctx).Model(&pmmodel.PointsProduct{}).
		Where("id = ?", id).
		Update("status", status).Error
}
