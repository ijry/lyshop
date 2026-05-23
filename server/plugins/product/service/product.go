package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
)

type ProductListQuery struct {
	CategoryID uint64 `form:"category_id"`
	Keyword    string `form:"keyword"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
}

type ProductDetail struct {
	productmodel.Product
	SKUs   []productmodel.ProductSku   `json:"skus"`
	Images []productmodel.ProductImage `json:"images"`
}

var defaultDetailJSON = json.RawMessage(`{"version":1,"blocks":[]}`)

func normalizeDetail(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 || string(raw) == "null" {
		return defaultDetailJSON
	}
	var payload productmodel.ProductDetailContent
	if err := json.Unmarshal(raw, &payload); err != nil {
		return defaultDetailJSON
	}
	if payload.Version <= 0 {
		payload.Version = 1
	}
	if payload.Blocks == nil {
		payload.Blocks = []productmodel.DetailBlock{}
	}
	normalized, err := json.Marshal(payload)
	if err != nil {
		return defaultDetailJSON
	}
	return normalized
}

func ListProducts(ctx context.Context, q ProductListQuery) ([]productmodel.Product, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 || q.Size > 100 {
		q.Size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&productmodel.Product{})
	if q.CategoryID > 0 {
		tx = tx.Where("category_id = ?", q.CategoryID)
	}
	if q.Keyword != "" {
		tx = tx.Where("title LIKE ?", "%"+q.Keyword+"%")
	}
	var total int64
	tx.Count(&total)
	var list []productmodel.Product
	err := tx.Order("sort desc, id desc").
		Offset((q.Page - 1) * q.Size).Limit(q.Size).Find(&list).Error
	return list, total, err
}

func GetProduct(ctx context.Context, id uint64) (*ProductDetail, error) {
	var p productmodel.Product
	if err := db.DB.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商品不存在")
		}
		return nil, err
	}
	p.Detail = normalizeDetail(p.Detail)
	detail := &ProductDetail{Product: p}
	db.DB.WithContext(ctx).Where("product_id = ?", id).Find(&detail.SKUs)
	db.DB.WithContext(ctx).Where("product_id = ?", id).Order("sort asc").Find(&detail.Images)
	return detail, nil
}

func CreateProduct(ctx context.Context, p *productmodel.Product, skus []productmodel.ProductSku, images []productmodel.ProductImage) error {
	p.Detail = normalizeDetail(p.Detail)
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		for i := range skus {
			skus[i].ProductID = p.ID
		}
		if len(skus) > 0 {
			if err := tx.Create(&skus).Error; err != nil {
				return err
			}
		}
		for i := range images {
			images[i].ProductID = p.ID
		}
		if len(images) > 0 {
			tx.Create(&images)
		}
		return nil
	})
}

func UpdateProduct(ctx context.Context, id uint64, updates map[string]any) error {
	if detail, ok := updates["detail"]; ok {
		raw, err := json.Marshal(detail)
		if err != nil {
			return err
		}
		updates["detail"] = normalizeDetail(raw)
	}
	return db.DB.WithContext(ctx).Model(&productmodel.Product{}).Where("id = ?", id).Updates(updates).Error
}

func DeleteProduct(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("product_id = ?", id).Delete(&productmodel.ProductSku{})
		tx.Where("product_id = ?", id).Delete(&productmodel.ProductImage{})
		return tx.Delete(&productmodel.Product{}, id).Error
	})
}

func ReplaceProductImages(ctx context.Context, productID uint64, images []productmodel.ProductImage) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productID).Delete(&productmodel.ProductImage{}).Error; err != nil {
			return err
		}
		for i := range images {
			images[i].ProductID = productID
		}
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
