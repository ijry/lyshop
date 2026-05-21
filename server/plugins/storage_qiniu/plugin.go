package storage_qiniu

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

type qiniuDriver struct{ Zone, Bucket, AccessKey, SecretKey, Domain string }

func (d *qiniuDriver) Name() string { return "qiniu" }
func (d *qiniuDriver) Upload(_ context.Context, fh *multipart.FileHeader) (*driverStorage.UploadResult, error) {
	if d.Bucket == "" { return nil, fmt.Errorf("七牛云未配置 Bucket") }
	// Production: use github.com/qiniu/go-sdk/v7
	return &driverStorage.UploadResult{Path: fh.Filename, URL: d.Domain + "/" + fh.Filename, Size: fh.Size}, nil
}
func (d *qiniuDriver) Delete(_ context.Context, _ string) error { return nil }
func (d *qiniuDriver) GetURL(path string) string { return d.Domain + "/" + path }

type storageQiniuPlugin struct{ meta plugin.Meta; driver *qiniuDriver }

func init() {
	var m plugin.Meta; json.Unmarshal(metaJSON, &m)
	plugin.Register(&storageQiniuPlugin{meta: m, driver: &qiniuDriver{}})
}

func (p *storageQiniuPlugin) Meta() plugin.Meta { return p.meta }
func (p *storageQiniuPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (p *storageQiniuPlugin) Migrate(_ *gorm.DB) error { return nil }
func (p *storageQiniuPlugin) Install() error {
	kv := func(k string) string {
		var c model.ConfigKV
		if db.DB.Where("plugin=? AND key=?","storage_qiniu",k).First(&c).Error==nil { return c.Value }
		return ""
	}
	p.driver.Zone = kv("zone"); p.driver.Bucket = kv("bucket")
	p.driver.AccessKey = kv("access_key"); p.driver.SecretKey = kv("secret_key")
	p.driver.Domain = kv("domain")
	driverStorage.Register(p.driver); return nil
}
func (p *storageQiniuPlugin) Uninstall() error { return nil }
