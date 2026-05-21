package storage_oss

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

type ossDriver struct{ Endpoint, Bucket, AccessKeyID, AccessKeySecret string }

func (d *ossDriver) Name() string { return "oss" }
func (d *ossDriver) Upload(_ context.Context, fh *multipart.FileHeader) (*driverStorage.UploadResult, error) {
	if d.Bucket == "" { return nil, fmt.Errorf("阿里云OSS未配置 Bucket") }
	// Production: use github.com/aliyun/aliyun-oss-go-sdk
	return &driverStorage.UploadResult{Path: fh.Filename, URL: "https://" + d.Bucket + "." + d.Endpoint + "/" + fh.Filename, Size: fh.Size}, nil
}
func (d *ossDriver) Delete(_ context.Context, _ string) error { return nil }
func (d *ossDriver) GetURL(path string) string { return "https://" + d.Bucket + "." + d.Endpoint + "/" + path }

type storageOssPlugin struct{ meta plugin.Meta; driver *ossDriver }

func init() {
	var m plugin.Meta; json.Unmarshal(metaJSON, &m)
	plugin.Register(&storageOssPlugin{meta: m, driver: &ossDriver{}})
}

func (p *storageOssPlugin) Meta() plugin.Meta { return p.meta }
func (p *storageOssPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (p *storageOssPlugin) Migrate(_ *gorm.DB) error { return nil }
func (p *storageOssPlugin) Install() error {
	kv := func(k string) string {
		var c model.ConfigKV
		if db.DB.Where("plugin=? AND key=?","storage_oss",k).First(&c).Error==nil { return c.Value }
		return ""
	}
	p.driver.Endpoint = kv("endpoint"); p.driver.Bucket = kv("bucket")
	p.driver.AccessKeyID = kv("access_key_id"); p.driver.AccessKeySecret = kv("access_key_secret")
	driverStorage.Register(p.driver); return nil
}
func (p *storageOssPlugin) Uninstall() error { return nil }
