package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ijry/lyshop/core/db"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
)

func ListWarehouses(ctx context.Context, q WarehouseListQuery) ([]wmsmodel.Warehouse, int64, error) {
	page, size := normalizePage(q.Page, q.Size)
	tx := db.DB.WithContext(ctx).Model(&wmsmodel.Warehouse{})
	if q.Keyword != "" {
		like := "%" + strings.TrimSpace(q.Keyword) + "%"
		tx = tx.Where("name LIKE ? OR code LIKE ? OR contact LIKE ? OR phone LIKE ?", like, like, like, like)
	}
	if q.Status != nil {
		tx = tx.Where("status = ?", *q.Status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []wmsmodel.Warehouse
	if err := tx.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func CreateWarehouse(ctx context.Context, w *wmsmodel.Warehouse) error {
	w.Name = strings.TrimSpace(w.Name)
	w.Code = strings.TrimSpace(w.Code)
	if w.Name == "" {
		return fmt.Errorf("仓库名称不能为空")
	}
	if w.Code == "" {
		return fmt.Errorf("仓库编码不能为空")
	}
	if w.Status != wmsmodel.WarehouseStatusEnabled && w.Status != wmsmodel.WarehouseStatusDisabled {
		w.Status = wmsmodel.WarehouseStatusEnabled
	}
	return db.DB.WithContext(ctx).Create(w).Error
}

func UpdateWarehouse(ctx context.Context, id uint64, in *wmsmodel.Warehouse) error {
	if id == 0 {
		return fmt.Errorf("仓库ID不能为空")
	}
	var row wmsmodel.Warehouse
	if err := db.DB.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return fmt.Errorf("仓库不存在")
	}
	name := strings.TrimSpace(in.Name)
	code := strings.TrimSpace(in.Code)
	if name == "" {
		return fmt.Errorf("仓库名称不能为空")
	}
	if code == "" {
		return fmt.Errorf("仓库编码不能为空")
	}
	updates := map[string]any{
		"code":    code,
		"name":    name,
		"address": strings.TrimSpace(in.Address),
		"contact": strings.TrimSpace(in.Contact),
		"phone":   strings.TrimSpace(in.Phone),
	}
	return db.DB.WithContext(ctx).Model(&wmsmodel.Warehouse{}).Where("id = ?", id).Updates(updates).Error
}

func UpdateWarehouseStatus(ctx context.Context, id uint64, status int8) error {
	if id == 0 {
		return fmt.Errorf("仓库ID不能为空")
	}
	if status != wmsmodel.WarehouseStatusEnabled && status != wmsmodel.WarehouseStatusDisabled {
		return fmt.Errorf("仓库状态非法")
	}
	tx := db.DB.WithContext(ctx).Model(&wmsmodel.Warehouse{}).Where("id = ?", id).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("仓库不存在")
	}
	return nil
}
