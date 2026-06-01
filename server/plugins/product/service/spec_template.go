package service

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"

	"github.com/ijry/lyshop/core/db"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"gorm.io/gorm"
)

type SpecTemplateListQuery struct {
	Keyword    string `form:"keyword"`
	CategoryID uint64 `form:"category_id"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
}

func normalizeCategoryIDs(ids []uint64) []uint64 {
	unique := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, exists := unique[id]; exists {
			continue
		}
		unique[id] = struct{}{}
		out = append(out, id)
	}
	slices.Sort(out)
	return out
}

func decodeCategoryIDs(raw json.RawMessage) []uint64 {
	if len(raw) == 0 || string(raw) == "null" {
		return []uint64{}
	}
	var ids []uint64
	if err := json.Unmarshal(raw, &ids); err != nil {
		return []uint64{}
	}
	return normalizeCategoryIDs(ids)
}

func BuildSpecTemplateCategoryIDs(ids []uint64) (json.RawMessage, error) {
	return json.Marshal(normalizeCategoryIDs(ids))
}

func decodeSpecTemplateAttrs(raw json.RawMessage) []SpecSchemaGroup {
	if len(raw) == 0 || string(raw) == "null" {
		return []SpecSchemaGroup{}
	}
	var attrs []SpecSchemaGroup
	if err := json.Unmarshal(raw, &attrs); err != nil {
		return []SpecSchemaGroup{}
	}
	return normalizeSpecSchema(attrs)
}

func BuildSpecTemplateAttrs(attrs []SpecSchemaGroup) (json.RawMessage, error) {
	return json.Marshal(normalizeSpecSchema(attrs))
}

func normalizeSpecTemplateModel(row *productmodel.SpecTemplate) {
	if row == nil {
		return
	}
	categoryIDs, err := BuildSpecTemplateCategoryIDs(decodeCategoryIDs(row.CategoryIDs))
	if err == nil {
		row.CategoryIDs = categoryIDs
	}
	attrs, err := BuildSpecTemplateAttrs(decodeSpecTemplateAttrs(row.Attrs))
	if err == nil {
		row.Attrs = attrs
	}
	if row.Status != 0 {
		row.Status = 1
	}
}

func ListSpecTemplates(ctx context.Context, q SpecTemplateListQuery) ([]productmodel.SpecTemplate, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 || q.Size > 200 {
		q.Size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&productmodel.SpecTemplate{})
	if strings.TrimSpace(q.Keyword) != "" {
		tx = tx.Where("name LIKE ?", "%"+strings.TrimSpace(q.Keyword)+"%")
	}

	var rows []productmodel.SpecTemplate
	if err := tx.Order("sort desc, id desc").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	filtered := make([]productmodel.SpecTemplate, 0, len(rows))
	for _, row := range rows {
		if q.CategoryID > 0 {
			ids := decodeCategoryIDs(row.CategoryIDs)
			if !slices.Contains(ids, q.CategoryID) {
				continue
			}
		}
		normalizeSpecTemplateModel(&row)
		filtered = append(filtered, row)
	}

	total := int64(len(filtered))
	offset := (q.Page - 1) * q.Size
	if offset >= len(filtered) {
		return []productmodel.SpecTemplate{}, total, nil
	}
	end := offset + q.Size
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func GetSpecTemplate(ctx context.Context, id uint64) (*productmodel.SpecTemplate, error) {
	var row productmodel.SpecTemplate
	if err := db.DB.WithContext(ctx).First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("规格模板不存在")
		}
		return nil, err
	}
	normalizeSpecTemplateModel(&row)
	return &row, nil
}

func CreateSpecTemplate(ctx context.Context, row *productmodel.SpecTemplate) error {
	normalizeSpecTemplateModel(row)
	return db.DB.WithContext(ctx).Create(row).Error
}

func UpdateSpecTemplate(ctx context.Context, id uint64, updates map[string]any) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row productmodel.SpecTemplate
		if err := tx.First(&row, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("规格模板不存在")
			}
			return err
		}

		if value, exists := updates["name"]; exists {
			if parsed, ok := value.(string); ok {
				row.Name = strings.TrimSpace(parsed)
			}
		}
		if value, exists := updates["status"]; exists {
			switch v := value.(type) {
			case int:
				if v == 0 {
					row.Status = 0
				} else {
					row.Status = 1
				}
			case int8:
				if v == 0 {
					row.Status = 0
				} else {
					row.Status = 1
				}
			case float64:
				if int(v) == 0 {
					row.Status = 0
				} else {
					row.Status = 1
				}
			}
		}
		if value, exists := updates["sort"]; exists {
			switch v := value.(type) {
			case int:
				row.Sort = v
			case float64:
				row.Sort = int(v)
			}
		}
		if value, exists := updates["category_ids"]; exists {
			if ids, ok := value.([]uint64); ok {
				raw, err := BuildSpecTemplateCategoryIDs(ids)
				if err != nil {
					return err
				}
				row.CategoryIDs = raw
			}
		}
		if value, exists := updates["attrs"]; exists {
			if attrs, ok := value.([]SpecSchemaGroup); ok {
				raw, err := BuildSpecTemplateAttrs(attrs)
				if err != nil {
					return err
				}
				row.Attrs = raw
			}
		}
		if strings.TrimSpace(row.Name) == "" {
			return errors.New("模板名称不能为空")
		}
		normalizeSpecTemplateModel(&row)
		return tx.Model(&productmodel.SpecTemplate{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"name":         row.Name,
				"category_ids": row.CategoryIDs,
				"attrs":        row.Attrs,
				"status":       row.Status,
				"sort":         row.Sort,
			}).Error
	})
}

func DeleteSpecTemplate(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Delete(&productmodel.SpecTemplate{}, id).Error
}
