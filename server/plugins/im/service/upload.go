package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	immodel "github.com/ijry/lyshop/plugins/im/model"
)

const maxIMUploadSize int64 = 10 << 20

type AttachmentInfo struct {
	URL         string `json:"url"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Mime        string `json:"mime"`
	MessageType string `json:"message_type"`
}

var imageExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
var fileExts = map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".txt": true, ".csv": true, ".md": true, ".zip": true}

func ClassifyAttachment(filename, mime string, size int64) (*AttachmentInfo, error) {
	if size > maxIMUploadSize {
		return nil, fmt.Errorf("文件过大，最大支持 10MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	lowerMime := strings.ToLower(strings.TrimSpace(mime))
	switch {
	case imageExts[ext] && strings.HasPrefix(lowerMime, "image/"):
		return &AttachmentInfo{Name: filepath.Base(filename), Size: size, Mime: mime, MessageType: immodel.MsgTypeImage}, nil
	case fileExts[ext]:
		return &AttachmentInfo{Name: filepath.Base(filename), Size: size, Mime: mime, MessageType: immodel.MsgTypeFile}, nil
	default:
		return nil, fmt.Errorf("不支持的文件类型：%s", ext)
	}
}

func UploadAttachment(ctx context.Context, fh *multipart.FileHeader) (*AttachmentInfo, error) {
	info, err := ClassifyAttachment(fh.Filename, fh.Header.Get("Content-Type"), fh.Size)
	if err != nil {
		return nil, err
	}
	driver, err := driverStorage.Get()
	if err != nil {
		return nil, err
	}
	res, err := driver.Upload(ctx, fh)
	if err != nil {
		return nil, err
	}
	info.URL = res.URL
	info.Path = res.Path
	if info.Mime == "" {
		info.Mime = res.Mime
	}
	if info.Size == 0 {
		info.Size = res.Size
	}
	return info, nil
}
