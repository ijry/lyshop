package storage_qiniu

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/model"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	qiniustorage "github.com/qiniu/go-sdk/v7/storage"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type qiniuDriver struct{ Zone, Bucket, AccessKey, SecretKey, Domain string }

func (d *qiniuDriver) Name() string { return "qiniu" }

func (d *qiniuDriver) Upload(ctx context.Context, fh *multipart.FileHeader) (*driverStorage.UploadResult, error) {
	if err := d.validate(); err != nil {
		return nil, err
	}
	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	cfg := d.uploadConfig()
	uploader := qiniustorage.NewFormUploader(cfg)
	putPolicy := qiniustorage.PutPolicy{Scope: strings.TrimSpace(d.Bucket)}
	mac := qbox.NewMac(strings.TrimSpace(d.AccessKey), strings.TrimSpace(d.SecretKey))
	upToken := putPolicy.UploadToken(mac)
	objectKey := buildObjectKey(fh.Filename)

	putExtra := qiniustorage.PutExtra{}
	if mime := strings.TrimSpace(fh.Header.Get("Content-Type")); mime != "" {
		putExtra.MimeType = mime
	}
	ret := qiniustorage.PutRet{}
	if err := uploader.Put(ctx, &ret, upToken, objectKey, src, fh.Size, &putExtra); err != nil {
		return nil, fmt.Errorf("上传到七牛云失败: %w", err)
	}
	return &driverStorage.UploadResult{
		Path: objectKey,
		URL:  d.GetURL(objectKey),
		Size: fh.Size,
		Mime: strings.TrimSpace(fh.Header.Get("Content-Type")),
	}, nil
}

func (d *qiniuDriver) Delete(_ context.Context, path string) error {
	if err := d.validate(); err != nil {
		return err
	}
	objectKey := cleanObjectPath(path)
	if objectKey == "" {
		return fmt.Errorf("七牛云删除路径不能为空")
	}
	mac := qbox.NewMac(strings.TrimSpace(d.AccessKey), strings.TrimSpace(d.SecretKey))
	bucketManager := qiniustorage.NewBucketManager(mac, d.uploadConfig())
	if err := bucketManager.Delete(strings.TrimSpace(d.Bucket), objectKey); err != nil {
		return fmt.Errorf("从七牛云删除文件失败: %w", err)
	}
	return nil
}

func (d *qiniuDriver) GetURL(path string) string {
	objectKey := cleanObjectPath(path)
	if objectKey == "" {
		return ""
	}
	return d.normalizedDomain() + "/" + objectKey
}

func (d *qiniuDriver) validate() error {
	if strings.TrimSpace(d.Bucket) == "" {
		return fmt.Errorf("七牛云未配置 Bucket")
	}
	if strings.TrimSpace(d.AccessKey) == "" {
		return fmt.Errorf("七牛云未配置 AccessKey")
	}
	if strings.TrimSpace(d.SecretKey) == "" {
		return fmt.Errorf("七牛云未配置 SecretKey")
	}
	if strings.TrimSpace(d.Domain) == "" {
		return fmt.Errorf("七牛云未配置访问域名")
	}
	return nil
}

func (d *qiniuDriver) normalizedDomain() string {
	domain := strings.TrimSpace(d.Domain)
	if domain == "" {
		return ""
	}
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}
	return strings.TrimSuffix(domain, "/")
}

func (d *qiniuDriver) uploadConfig() *qiniustorage.Config {
	cfg := &qiniustorage.Config{
		UseHTTPS: true,
	}
	if zone := mapZone(strings.TrimSpace(d.Zone)); zone != nil {
		cfg.Region = zone
	}
	return cfg
}

func mapZone(zone string) *qiniustorage.Region {
	switch strings.ToLower(strings.TrimSpace(zone)) {
	case "z0", "huadong", "cn-east-1":
		region := qiniustorage.ZoneHuadong
		return &region
	case "z1", "huabei", "cn-north-1":
		region := qiniustorage.ZoneHuabei
		return &region
	case "z2", "huanan", "cn-south-1":
		region := qiniustorage.ZoneHuanan
		return &region
	case "na0", "beimei", "us-north-1":
		region := qiniustorage.ZoneBeimei
		return &region
	case "as0", "xinjiapo", "ap-southeast-1":
		region := qiniustorage.ZoneXinjiapo
		return &region
	default:
		return nil
	}
}

func buildObjectKey(filename string) string {
	ext := filepath.Ext(filename)
	now := time.Now()
	suffix := randomHex(6)
	if suffix == "" {
		return fmt.Sprintf("%s/%d%s", now.Format("2006/01"), now.UnixNano(), ext)
	}
	return fmt.Sprintf("%s/%d_%s%s", now.Format("2006/01"), now.UnixNano(), suffix, ext)
}

func randomHex(size int) string {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}

func cleanObjectPath(path string) string {
	objectPath := strings.TrimSpace(path)
	if objectPath == "" {
		return ""
	}
	if u, err := url.Parse(objectPath); err == nil && u.Host != "" {
		objectPath = u.Path
	}
	return strings.TrimPrefix(objectPath, "/")
}

type storageQiniuPlugin struct {
	meta   plugin.Meta
	driver *qiniuDriver
}

func init() {
	var m plugin.Meta
	json.Unmarshal(metaJSON, &m)
	plugin.Register(&storageQiniuPlugin{meta: m, driver: &qiniuDriver{}})
}

func (p *storageQiniuPlugin) Meta() plugin.Meta                    { return p.meta }
func (p *storageQiniuPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (p *storageQiniuPlugin) Migrate(_ *gorm.DB) error             { return nil }
func (p *storageQiniuPlugin) Install() error {
	kv := func(k string) string {
		var c model.ConfigKV
		if db.DB.Where("plugin=? AND key=?", "storage_qiniu", k).First(&c).Error == nil {
			return c.Value
		}
		return ""
	}
	p.driver.Zone = kv("zone")
	p.driver.Bucket = kv("bucket")
	p.driver.AccessKey = kv("access_key")
	p.driver.SecretKey = kv("secret_key")
	p.driver.Domain = kv("domain")
	driverStorage.Register(p.driver)
	return nil
}
func (p *storageQiniuPlugin) Uninstall() error { return nil }
