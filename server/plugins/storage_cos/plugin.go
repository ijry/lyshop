package storage_cos

import (
	_ "embed"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type cosDriver struct{ Region, Bucket, SecretID, SecretKey string }

func (d *cosDriver) Name() string { return "cos" }
func (d *cosDriver) Upload(_ context.Context, fh *multipart.FileHeader) (*driverStorage.UploadResult, error) {
	if d.Bucket == "" { return nil, fmt.Errorf("腾讯云COS未配置 Bucket") }
	// Production: use github.com/tencentyun/cos-go-sdk-v5
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", d.Bucket, d.Region, fh.Filename)
	return &driverStorage.UploadResult{Path: fh.Filename, URL: url, Size: fh.Size}, nil
}
func (d *cosDriver) Delete(_ context.Context, _ string) error { return nil }
func (d *cosDriver) GetURL(path string) string {
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", d.Bucket, d.Region, path)
}

type storageCosPlugin struct{ meta plugin.Meta; driver *cosDriver }

func init() {
	var m plugin.Meta; json.Unmarshal(metaJSON, &m)
	plugin.Register(&storageCosPlugin{meta: m, driver: &cosDriver{}})
}

func (p *storageCosPlugin) Meta() plugin.Meta { return p.meta }
func (p *storageCosPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (p *storageCosPlugin) Migrate(_ *gorm.DB) error { return nil }
func (p *storageCosPlugin) Install() error {
	kv := func(k string) string {
		var c model.ConfigKV
		if db.DB.Where("plugin=? AND key=?","storage_cos",k).First(&c).Error==nil { return c.Value }
		return ""
	}
	p.driver.Region = kv("region"); p.driver.Bucket = kv("bucket")
	p.driver.SecretID = kv("secret_id"); p.driver.SecretKey = kv("secret_key")
	driverStorage.Register(p.driver); return nil
}
func (p *storageCosPlugin) Uninstall() error { return nil }
