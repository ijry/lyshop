package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"
)

// UploadResult holds the stored path and public URL.
type UploadResult struct {
	Path string
	URL  string
	Size int64
	Mime string
}

// Driver is the interface all file storage plugins must implement.
type Driver interface {
	Name() string
	Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
	Delete(ctx context.Context, path string) error
	GetURL(path string) string
}

var (
	mu     sync.RWMutex
	active Driver
)

func Register(d Driver) { mu.Lock(); active = d; mu.Unlock() }

func Get() (Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	if active == nil {
		return nil, fmt.Errorf("no storage driver registered")
	}
	return active, nil
}
