package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"sync"
)

// UploadResult holds the stored path and public URL.
type UploadResult struct {
	Path string `json:"path"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Mime string `json:"mime"`
}

// Driver is the interface all file storage plugins must implement.
type Driver interface {
	Name() string
	Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
	Delete(ctx context.Context, path string) error
	GetURL(path string) string
}

var (
	mu          sync.RWMutex
	active      Driver
	drivers     = map[string]Driver{}
	defaultName string
)

func Register(d Driver) {
	mu.Lock()
	defer mu.Unlock()
	if d == nil {
		return
	}
	key := normalizeName(d.Name())
	drivers[key] = d
	active = d
	if defaultName == "" {
		defaultName = key
	}
}

func Get() (Driver, error) {
	return GetByName("")
}

func SetDefault(name string) {
	mu.Lock()
	defer mu.Unlock()
	defaultName = normalizeName(name)
}

func GetDefaultName() string {
	mu.RLock()
	defer mu.RUnlock()
	return defaultName
}

func GetByName(name string) (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	requested := normalizeName(name)
	if requested != "" {
		if d, ok := drivers[requested]; ok {
			return d, nil
		}
		return nil, fmt.Errorf("storage driver %q not registered", requested)
	}
	if defaultName != "" {
		if d, ok := drivers[defaultName]; ok {
			return d, nil
		}
	}
	if active != nil {
		return active, nil
	}
	return nil, fmt.Errorf("no storage driver registered")
}

func normalizeName(name string) string {
	v := strings.ToLower(strings.TrimSpace(name))
	switch v {
	case "storage_local", "local":
		return "local"
	case "storage_oss", "oss", "aliyun_oss":
		return "oss"
	case "storage_cos", "cos", "qcloud_cos":
		return "cos"
	case "storage_qiniu", "qiniu", "qiniu_kodo":
		return "qiniu"
	default:
		return v
	}
}
