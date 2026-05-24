package service

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
)

func ListCategories(ctx context.Context, includeDisabled bool) ([]productmodel.ProductCategory, error) {
	var list []productmodel.ProductCategory
	tx := db.DB.WithContext(ctx).Model(&productmodel.ProductCategory{})
	if !includeDisabled {
		tx = tx.Where("status = 1")
	}
	err := tx.Order("sort asc, id asc").Find(&list).Error
	return list, err
}

func CreateCategory(ctx context.Context, c *productmodel.ProductCategory) error {
	return db.DB.WithContext(ctx).Create(c).Error
}

func UpdateCategory(ctx context.Context, id uint64, updates map[string]any) error {
	return db.DB.WithContext(ctx).Model(&productmodel.ProductCategory{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteCategory(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Delete(&productmodel.ProductCategory{}, id).Error
}
