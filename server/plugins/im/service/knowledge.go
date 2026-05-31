package service

import (
	"context"

	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
)

// ListKnowledge returns paginated knowledge-base entries (optionally filtered
// by a keyword over title/content/tags).
func ListKnowledge(ctx context.Context, keyword string, page, size int) ([]immodel.ImKnowledge, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&immodel.ImKnowledge{})
	if keyword != "" {
		like := "%" + keyword + "%"
		tx = tx.Where("title LIKE ? OR content LIKE ? OR tags LIKE ?", like, like, like)
	}
	var total int64
	tx.Count(&total)
	var list []immodel.ImKnowledge
	err := tx.Order("sort asc, id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

// CreateKnowledge persists a new entry and best-effort embeds it.
func CreateKnowledge(ctx context.Context, k *immodel.ImKnowledge) error {
	if err := db.DB.WithContext(ctx).Create(k).Error; err != nil {
		return err
	}
	go EmbedKnowledgeEntry(context.Background(), k.ID)
	return nil
}

// UpdateKnowledge updates an entry's editable fields and re-embeds it.
func UpdateKnowledge(ctx context.Context, id uint64, fields map[string]any) error {
	// Content changed → mark stale until re-embedded.
	fields["indexed"] = 0
	fields["embedding"] = nil
	if err := db.DB.WithContext(ctx).Model(&immodel.ImKnowledge{}).
		Where("id = ?", id).Updates(fields).Error; err != nil {
		return err
	}
	go EmbedKnowledgeEntry(context.Background(), id)
	return nil
}

// DeleteKnowledge removes an entry.
func DeleteKnowledge(ctx context.Context, id uint64) error {
	return db.DB.WithContext(ctx).Delete(&immodel.ImKnowledge{}, id).Error
}
