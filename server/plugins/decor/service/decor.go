package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ijry/lyshop/core/db"
	decormodel "github.com/ijry/lyshop/plugins/decor/model"
	"gorm.io/gorm"
)

// GetPage returns the current page config for pageKey (merchantID=0 for single-tenant).
func GetPage(ctx context.Context, merchantID uint64, pageKey string) (*decormodel.DecorPage, error) {
	var page decormodel.DecorPage
	err := db.DB.WithContext(ctx).
		Where("merchant_id = ? AND page_key = ?", merchantID, pageKey).
		First(&page).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Return empty page if not yet configured
		empty, _ := json.Marshal([]any{})
		return &decormodel.DecorPage{
			MerchantID: merchantID,
			PageKey:    pageKey,
			Components: empty,
		}, nil
	}
	return &page, err
}

// SavePage upserts the page configuration.
func SavePage(ctx context.Context, merchantID uint64, pageKey string, components json.RawMessage) (*decormodel.DecorPage, error) {
	var page decormodel.DecorPage
	err := db.DB.WithContext(ctx).
		Where("merchant_id = ? AND page_key = ?", merchantID, pageKey).
		First(&page).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		page = decormodel.DecorPage{
			MerchantID: merchantID,
			PageKey:    pageKey,
			Components: components,
		}
		return &page, db.DB.WithContext(ctx).Create(&page).Error
	}
	page.Components = components
	return &page, db.DB.WithContext(ctx).Save(&page).Error
}

// PublishPage marks the page as published.
func PublishPage(ctx context.Context, merchantID uint64, pageKey string) error {
	now := time.Now()
	return db.DB.WithContext(ctx).Model(&decormodel.DecorPage{}).
		Where("merchant_id = ? AND page_key = ?", merchantID, pageKey).
		Update("published_at", now).Error
}
