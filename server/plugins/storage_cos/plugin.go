package storage_cos

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/model"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type cosDriver struct{ Region, Bucket, SecretID, SecretKey string }

func (d *cosDriver) Name() string { return "cos" }
func (d *cosDriver) Upload(ctx context.Context, fh *multipart.FileHeader) (*driverStorage.UploadResult, error) {
	if err := d.validate(); err != nil {
		return nil, err
	}
	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	client, err := d.client()
	if err != nil {
		return nil, err
	}
	objectKey := buildObjectKey(fh.Filename)
	opt := &cos.ObjectPutOptions{}
	if mime := strings.TrimSpace(fh.Header.Get("Content-Type")); mime != "" {
		opt.ObjectPutHeaderOptions = &cos.ObjectPutHeaderOptions{
			ContentType: mime,
		}
	}
	if _, err := client.Object.Put(ctx, objectKey, src, opt); err != nil {
		return nil, fmt.Errorf("上传到腾讯云COS失败: %w", err)
	}
	return &driverStorage.UploadResult{
		Path: objectKey,
		URL:  d.GetURL(objectKey),
		Size: fh.Size,
		Mime: strings.TrimSpace(fh.Header.Get("Content-Type")),
	}, nil
}

func (d *cosDriver) Delete(ctx context.Context, path string) error {
	if err := d.validate(); err != nil {
		return err
	}
	objectKey := cleanObjectPath(path)
	if objectKey == "" {
		return fmt.Errorf("腾讯云COS删除路径不能为空")
	}
	client, err := d.client()
	if err != nil {
		return err
	}
	resp, err := client.Object.Delete(ctx, objectKey)
	if err != nil {
		return fmt.Errorf("从腾讯云COS删除文件失败: %w", err)
	}
	if resp != nil {
		_ = resp.Body.Close()
	}
	return nil
}

func (d *cosDriver) GetURL(path string) string {
	objectKey := cleanObjectPath(path)
	if objectKey == "" {
		return ""
	}
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", d.normalizedBucket(), d.normalizedRegion(), objectKey)
}

func (d *cosDriver) validate() error {
	if strings.TrimSpace(d.Region) == "" {
		return fmt.Errorf("腾讯云COS未配置 Region")
	}
	if strings.TrimSpace(d.Bucket) == "" {
		return fmt.Errorf("腾讯云COS未配置 Bucket")
	}
	if strings.TrimSpace(d.SecretID) == "" {
		return fmt.Errorf("腾讯云COS未配置 SecretId")
	}
	if strings.TrimSpace(d.SecretKey) == "" {
		return fmt.Errorf("腾讯云COS未配置 SecretKey")
	}
	return nil
}

func (d *cosDriver) normalizedBucket() string {
	bucket := strings.TrimSpace(d.Bucket)
	if idx := strings.LastIndex(bucket, "-"); idx > 0 {
		bucketName := bucket[:idx]
		appid := bucket[idx+1:]
		if appid != "" {
			return bucket
		}
		return bucketName
	}
	return bucket
}

func (d *cosDriver) normalizedRegion() string {
	return strings.TrimSpace(d.Region)
}

func (d *cosDriver) client() (*cos.Client, error) {
	baseURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", d.normalizedBucket(), d.normalizedRegion())
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("解析腾讯云COS地址失败: %w", err)
	}
	b := &cos.BaseURL{BucketURL: u}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  strings.TrimSpace(d.SecretID),
			SecretKey: strings.TrimSpace(d.SecretKey),
		},
	}), nil
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

type storageCosPlugin struct {
	meta   plugin.Meta
	driver *cosDriver
}

func init() {
	var m plugin.Meta
	json.Unmarshal(metaJSON, &m)
	plugin.Register(&storageCosPlugin{meta: m, driver: &cosDriver{}})
}

func (p *storageCosPlugin) Meta() plugin.Meta                    { return p.meta }
func (p *storageCosPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (p *storageCosPlugin) Migrate(_ *gorm.DB) error             { return nil }
func (p *storageCosPlugin) Install() error {
	kv := func(k string) string {
		var c model.ConfigKV
		if db.DB.Where("plugin=? AND key=?", "storage_cos", k).First(&c).Error == nil {
			return c.Value
		}
		return ""
	}
	p.driver.Region = kv("region")
	p.driver.Bucket = kv("bucket")
	p.driver.SecretID = kv("secret_id")
	p.driver.SecretKey = kv("secret_key")
	driverStorage.Register(p.driver)
	return nil
}
func (p *storageCosPlugin) Uninstall() error { return nil }
