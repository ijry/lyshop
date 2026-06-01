package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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

// ImportResult summarises a document import.
type ImportResult struct {
	Filename string `json:"filename"`
	Chunks   int    `json:"chunks"` // number of knowledge entries created
}

// ImportDocument extracts text from an uploaded enterprise document, slices it
// into overlapping chunks, and stores each chunk as an ImKnowledge entry.
//
// title prefixes each chunk's title (defaults to the filename); tags are shared
// across all chunks. chunkSize/overlap are measured in runes (0 → defaults).
// Created entries are embedded in the background like CreateKnowledge.
func ImportDocument(ctx context.Context, filename string, data []byte, title, tags string, chunkSize, overlap int) (*ImportResult, error) {
	if !IsSupportedDoc(filename) {
		return nil, fmt.Errorf("不支持的文件格式：%s（支持 %s）", filepath.Ext(filename), strings.Join(SupportedDocExts, " "))
	}
	text, err := ExtractText(filename, data)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("未能从文档中提取到文本")
	}

	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	if strings.TrimSpace(title) == "" {
		title = base
	}

	chunks := SliceText(text, chunkSize, overlap)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("切片结果为空")
	}

	res := &ImportResult{Filename: filepath.Base(filename)}
	for i, chunk := range chunks {
		k := &immodel.ImKnowledge{
			Title:   fmt.Sprintf("%s (%d/%d)", title, i+1, len(chunks)),
			Content: chunk,
			Tags:    tags,
			Status:  1,
			Sort:    i,
		}
		if err := db.DB.WithContext(ctx).Create(k).Error; err != nil {
			return res, err
		}
		res.Chunks++
		id := k.ID
		go EmbedKnowledgeEntry(context.Background(), id)
	}
	return res, nil
}
