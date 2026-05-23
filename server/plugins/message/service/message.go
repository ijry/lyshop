package service

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	msgmodel "github.com/ijry/lyshop/plugins/message/model"
)

func Send(ctx context.Context, userID uint64, group, title, content string) error {
	return db.DB.WithContext(ctx).Create(&msgmodel.Message{
		UserID: userID, Group: group, Title: title, Content: content,
	}).Error
}

func ListByGroup(ctx context.Context, userID uint64, group string, page, size int) ([]msgmodel.Message, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 50 { size = 20 }
	tx := db.DB.WithContext(ctx).Where("(user_id = ? OR user_id = 0)", userID)
	if group != "" { tx = tx.Where("`group` = ?", group) }
	var total int64
	tx.Model(&msgmodel.Message{}).Count(&total)
	var list []msgmodel.Message
	err := tx.Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error
	return list, total, err
}

func MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	return db.DB.WithContext(ctx).Model(&msgmodel.Message{}).
		Where("id IN ? AND (user_id = ? OR user_id = 0)", ids, userID).
		Update("is_read", 1).Error
}

func UnreadCounts(ctx context.Context, userID uint64) (map[string]int64, error) {
	type result struct {
		Group string
		Count int64
	}
	var results []result
	err := db.DB.WithContext(ctx).Model(&msgmodel.Message{}).
		Select("`group`, count(*) as count").
		Where("(user_id = ? OR user_id = 0) AND is_read = 0", userID).
		Group("`group`").Find(&results).Error
	m := map[string]int64{}
	for _, r := range results { m[r.Group] = r.Count }
	return m, err
}

func AdminList(ctx context.Context, group string, page, size int) ([]msgmodel.Message, int64, error) {
	if page <= 0 { page = 1 }
	if size <= 0 || size > 100 { size = 20 }
	tx := db.DB.WithContext(ctx).Model(&msgmodel.Message{})
	if group != "" { tx = tx.Where("`group` = ?", group) }
	var total int64
	tx.Count(&total)
	var list []msgmodel.Message
	err := tx.Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error
	return list, total, err
}
