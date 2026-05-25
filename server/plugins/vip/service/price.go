package service

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	vipmodel "github.com/ijry/lyshop/plugins/vip/model"
)

func GetVipSkuPrice(ctx context.Context, levelID, productID, skuID uint64) (*vipmodel.SkuPrice, error) {
	var row vipmodel.SkuPrice
	err := db.DB.WithContext(ctx).
		Where("status = 1 AND level_id = ? AND product_id = ? AND (sku_id = ? OR sku_id = 0)", levelID, productID, skuID).
		Order("sku_id desc, id desc").
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}
