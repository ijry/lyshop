package service

import (
	"context"
	"errors"
	"time"

	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FavoriteProductItem struct {
	productmodel.Product
	IsFavorited bool      `json:"is_favorited"`
	FavoritedAt time.Time `json:"favorited_at"`
}

func FavoriteProduct(ctx context.Context, userID, productID uint64) error {
	if userID == 0 || productID == 0 {
		return errors.New("参数错误")
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ensureProductExistsTx(tx, productID); err != nil {
			return err
		}
		fav := productmodel.ProductFavorite{UserID: userID, ProductID: productID}
		result := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "product_id"}},
			DoNothing: true,
		}).Create(&fav)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return tx.Model(&productmodel.Product{}).
				Where("id = ?", productID).
				UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
		}
		return nil
	})
}

func UnfavoriteProduct(ctx context.Context, userID, productID uint64) error {
	if userID == 0 || productID == 0 {
		return errors.New("参数错误")
	}
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ? AND product_id = ?", userID, productID).
			Delete(&productmodel.ProductFavorite{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return tx.Model(&productmodel.Product{}).Where("id = ?", productID).
				UpdateColumn("favorite_count", gorm.Expr("CASE WHEN favorite_count > 0 THEN favorite_count - 1 ELSE 0 END")).Error
		}
		return nil
	})
}

func ListUserFavorites(ctx context.Context, userID uint64, page, size int) ([]FavoriteProductItem, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	tx := db.DB.WithContext(ctx).
		Table("product_favorites pf").
		Joins("JOIN products p ON p.id = pf.product_id")

	var total int64
	if err := tx.Where("pf.user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []FavoriteProductItem
	err := tx.Where("pf.user_id = ?", userID).
		Select("p.*, pf.created_at AS favorited_at").
		Order("pf.created_at DESC, pf.id DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	for i := range rows {
		rows[i].IsFavorited = true
		rows[i].Detail = normalizeDetail(rows[i].Detail)
	}
	return rows, total, nil
}

func getFavoritedProductIDSet(ctx context.Context, userID uint64, productIDs []uint64) (map[uint64]struct{}, error) {
	set := make(map[uint64]struct{})
	if userID == 0 || len(productIDs) == 0 {
		return set, nil
	}
	var rows []productmodel.ProductFavorite
	err := db.DB.WithContext(ctx).
		Where("user_id = ? AND product_id IN ?", userID, productIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		set[row.ProductID] = struct{}{}
	}
	return set, nil
}

func ensureProductExistsTx(tx *gorm.DB, productID uint64) error {
	var count int64
	if err := tx.Model(&productmodel.Product{}).Where("id = ?", productID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("商品不存在")
	}
	return nil
}
