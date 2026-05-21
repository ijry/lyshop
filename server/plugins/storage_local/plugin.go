package storage_local

import (
	_ "embed"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

// localDriver stores files on the local filesystem.
type localDriver struct {
	uploadDir string
	baseURL   string
}

func (d *localDriver) Name() string { return "local" }

func (d *localDriver) Upload(_ context.Context, fh *multipart.FileHeader) (*driverStorage.UploadResult, error) {
	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := filepath.Ext(fh.Filename)
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := filepath.Join(d.uploadDir, time.Now().Format("2006/01"))
	os.MkdirAll(dir, 0755)
	dst := filepath.Join(dir, name)

	out, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	size, err := io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	relPath := filepath.ToSlash(filepath.Join(time.Now().Format("2006/01"), name))
	return &driverStorage.UploadResult{
		Path: relPath,
		URL:  d.baseURL + "/" + relPath,
		Size: size,
		Mime: fh.Header.Get("Content-Type"),
	}, nil
}

func (d *localDriver) Delete(_ context.Context, path string) error {
	return os.Remove(filepath.Join(d.uploadDir, path))
}

func (d *localDriver) GetURL(path string) string { return d.baseURL + "/" + path }

// Plugin

type storageLocalPlugin struct {
	meta   plugin.Meta
	driver *localDriver
}

func init() {
	var m plugin.Meta
	json.Unmarshal(metaJSON, &m)
	plugin.Register(&storageLocalPlugin{meta: m, driver: &localDriver{
		uploadDir: "uploads",
		baseURL:   "http://localhost:8080/uploads",
	}})
}

func (p *storageLocalPlugin) Meta() plugin.Meta { return p.meta }
func (p *storageLocalPlugin) RegisterRoutes(front, _ *gin.RouterGroup) {
	// Serve uploaded files statically — handled by Nginx in production
	// In dev: use gin.Static
	front.Static("/uploads", p.driver.uploadDir)
}
func (p *storageLocalPlugin) Migrate(_ *gorm.DB) error { return nil }
func (p *storageLocalPlugin) Install() error {
	loadKV := func(key, def string) string {
		var cfg model.ConfigKV
		if db.DB.Where("plugin = ? AND key = ?", "storage_local", key).First(&cfg).Error == nil {
			return cfg.Value
		}
		return def
	}
	p.driver.uploadDir = loadKV("upload_dir", "uploads")
	p.driver.baseURL = loadKV("base_url", "http://localhost:8080/uploads")
	os.MkdirAll(p.driver.uploadDir, 0755)
	driverStorage.Register(p.driver)
	return nil
}
func (p *storageLocalPlugin) Uninstall() error { return nil }
