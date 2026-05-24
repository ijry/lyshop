package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	decormodel "github.com/ijry/lyshop/plugins/decor/model"
	"gorm.io/gorm"
)

const (
	DefaultVariantKey  = "default"
	DefaultVariantName = "默认副本"
)

func normalizeVariantKey(variantKey string) string {
	key := strings.ToLower(strings.TrimSpace(variantKey))
	if key == "" {
		return DefaultVariantKey
	}
	key = strings.ReplaceAll(key, " ", "_")
	return key
}

func normalizeVariantName(variantName string) string {
	name := strings.TrimSpace(variantName)
	if name == "" {
		return DefaultVariantName
	}
	return name
}

func emptyComponents() json.RawMessage {
	empty, _ := json.Marshal([]any{})
	return empty
}

// GetPage returns the current page config for pageKey (merchantID=0 for single-tenant).
func GetPage(ctx context.Context, merchantID uint64, pageKey string, variantKey ...string) (*decormodel.DecorPage, error) {
	vKey := DefaultVariantKey
	if len(variantKey) > 0 {
		vKey = normalizeVariantKey(variantKey[0])
	}
	var page decormodel.DecorPage
	err := db.DB.WithContext(ctx).
		Where("merchant_id = ? AND page_key = ? AND variant_key = ?", merchantID, pageKey, vKey).
		First(&page).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Return empty page if not yet configured
		return &decormodel.DecorPage{
			MerchantID:  merchantID,
			PageKey:     pageKey,
			VariantKey:  vKey,
			VariantName: DefaultVariantName,
			Components:  emptyComponents(),
		}, nil
	}
	return &page, err
}

// SavePage upserts the page configuration.
func SavePage(ctx context.Context, merchantID uint64, pageKey string, components json.RawMessage, variantKey ...string) (*decormodel.DecorPage, error) {
	vKey := DefaultVariantKey
	if len(variantKey) > 0 {
		vKey = normalizeVariantKey(variantKey[0])
	}
	var page decormodel.DecorPage
	err := db.DB.WithContext(ctx).
		Where("merchant_id = ? AND page_key = ? AND variant_key = ?", merchantID, pageKey, vKey).
		First(&page).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		page = decormodel.DecorPage{
			MerchantID:  merchantID,
			PageKey:     pageKey,
			VariantKey:  vKey,
			VariantName: DefaultVariantName,
			Components:  components,
		}
		return &page, db.DB.WithContext(ctx).Create(&page).Error
	}
	page.Components = components
	return &page, db.DB.WithContext(ctx).Save(&page).Error
}

// PublishPage marks the page as published.
func PublishPage(ctx context.Context, merchantID uint64, pageKey string, variantKey ...string) error {
	vKey := DefaultVariantKey
	if len(variantKey) > 0 {
		vKey = normalizeVariantKey(variantKey[0])
	}
	page, err := GetPage(ctx, merchantID, pageKey, vKey)
	if err != nil {
		return err
	}
	if page.ID == 0 {
		created, saveErr := SavePage(ctx, merchantID, pageKey, page.Components, vKey)
		if saveErr != nil {
			return saveErr
		}
		page = created
	}
	now := time.Now()
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&decormodel.DecorPage{}).
			Where("merchant_id = ? AND page_key = ?", merchantID, pageKey).
			Updates(map[string]any{"is_published": false, "published_at": nil}).Error; err != nil {
			return err
		}
		return tx.Model(&decormodel.DecorPage{}).
			Where("id = ?", page.ID).
			Updates(map[string]any{"is_published": true, "published_at": now}).Error
	})
}

func GetPublishedPage(ctx context.Context, merchantID uint64, pageKey string) (*decormodel.DecorPage, error) {
	var page decormodel.DecorPage
	err := db.DB.WithContext(ctx).
		Where("merchant_id = ? AND page_key = ? AND is_published = ?", merchantID, pageKey, true).
		Order("published_at desc, id desc").
		First(&page).Error
	if err == nil {
		return &page, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return GetPage(ctx, merchantID, pageKey, DefaultVariantKey)
}

func ListVariants(ctx context.Context, merchantID uint64, pageKey string) ([]decormodel.DecorPage, error) {
	var rows []decormodel.DecorPage
	if err := db.DB.WithContext(ctx).
		Where("merchant_id = ? AND page_key = ?", merchantID, pageKey).
		Order("is_published desc, updated_at desc, id desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []decormodel.DecorPage{{
			MerchantID:  merchantID,
			PageKey:     pageKey,
			VariantKey:  DefaultVariantKey,
			VariantName: DefaultVariantName,
			Components:  emptyComponents(),
		}}, nil
	}
	return rows, nil
}

func CreateVariantCopy(ctx context.Context, merchantID uint64, pageKey string, fromVariantKey string, newVariantKey string, newVariantName string) (*decormodel.DecorPage, error) {
	src, err := GetPage(ctx, merchantID, pageKey, fromVariantKey)
	if err != nil {
		return nil, err
	}
	if src.ID == 0 {
		src.Components = emptyComponents()
	}
	targetKey := normalizeVariantKey(newVariantKey)
	targetName := normalizeVariantName(newVariantName)
	if targetKey == "" {
		return nil, fmt.Errorf("副本标识不能为空")
	}
	if targetKey == normalizeVariantKey(fromVariantKey) {
		return nil, fmt.Errorf("新副本标识不能与来源副本相同")
	}
	var existing decormodel.DecorPage
	if err := db.DB.WithContext(ctx).
		Where("merchant_id = ? AND page_key = ? AND variant_key = ?", merchantID, pageKey, targetKey).
		First(&existing).Error; err == nil {
		return nil, fmt.Errorf("副本标识已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	newPage := decormodel.DecorPage{
		MerchantID:  merchantID,
		PageKey:     pageKey,
		VariantKey:  targetKey,
		VariantName: targetName,
		Components:  src.Components,
		IsPublished: false,
		PublishedAt: nil,
	}
	return &newPage, db.DB.WithContext(ctx).Create(&newPage).Error
}

func RenameVariant(ctx context.Context, merchantID uint64, pageKey string, variantKey string, variantName string) error {
	vKey := normalizeVariantKey(variantKey)
	if vKey == "" {
		return fmt.Errorf("副本标识不能为空")
	}
	name := normalizeVariantName(variantName)
	return db.DB.WithContext(ctx).Model(&decormodel.DecorPage{}).
		Where("merchant_id = ? AND page_key = ? AND variant_key = ?", merchantID, pageKey, vKey).
		Update("variant_name", name).Error
}

func DeleteVariant(ctx context.Context, merchantID uint64, pageKey string, variantKey string) error {
	vKey := normalizeVariantKey(variantKey)
	if vKey == DefaultVariantKey {
		return fmt.Errorf("默认副本不支持删除")
	}
	var page decormodel.DecorPage
	if err := db.DB.WithContext(ctx).
		Where("merchant_id = ? AND page_key = ? AND variant_key = ?", merchantID, pageKey, vKey).
		First(&page).Error; err != nil {
		return err
	}
	if page.IsPublished {
		return fmt.Errorf("已发布副本不支持删除，请先发布其他副本")
	}
	return db.DB.WithContext(ctx).Delete(&page).Error
}
