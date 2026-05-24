package storage_oss

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

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
	if err := d.validate(); err != nil {
		return nil, err
	}
	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	bucket, err := d.bucketClient()
	if err != nil {
		return nil, err
	}
	objectKey := buildObjectKey(fh.Filename)
	options := []oss.Option{}
	if mime := strings.TrimSpace(fh.Header.Get("Content-Type")); mime != "" {
		options = append(options, oss.ContentType(mime))
	}
	if err := bucket.PutObject(objectKey, src, options...); err != nil {
		return nil, fmt.Errorf("上传到阿里云OSS失败: %w", err)
	}
	return &driverStorage.UploadResult{
		Path: objectKey,
		URL:  d.GetURL(objectKey),
		Size: fh.Size,
		Mime: strings.TrimSpace(fh.Header.Get("Content-Type")),
	}, nil
}

func (d *ossDriver) Delete(_ context.Context, path string) error {
	if err := d.validate(); err != nil {
		return err
	}
	objectKey := cleanObjectPath(path)
	if objectKey == "" {
		return fmt.Errorf("阿里云OSS删除路径不能为空")
	}
	bucket, err := d.bucketClient()
	if err != nil {
		return err
	}
	if err := bucket.DeleteObject(objectKey); err != nil {
		return fmt.Errorf("从阿里云OSS删除文件失败: %w", err)
	}
	return nil
}

func (d *ossDriver) GetURL(path string) string {
	objectKey := cleanObjectPath(path)
	if objectKey == "" {
		return ""
	}
	return fmt.Sprintf("https://%s.%s/%s", d.Bucket, d.normalizedEndpoint(), objectKey)
}

func (d *ossDriver) validate() error {
	if strings.TrimSpace(d.Endpoint) == "" {
		return fmt.Errorf("阿里云OSS未配置 Endpoint")
	}
	if strings.TrimSpace(d.Bucket) == "" {
		return fmt.Errorf("阿里云OSS未配置 Bucket")
	}
	if strings.TrimSpace(d.AccessKeyID) == "" {
		return fmt.Errorf("阿里云OSS未配置 AccessKeyId")
	}
	if strings.TrimSpace(d.AccessKeySecret) == "" {
		return fmt.Errorf("阿里云OSS未配置 AccessKeySecret")
	}
	return nil
}

func (d *ossDriver) normalizedEndpoint() string {
	endpoint := strings.TrimSpace(d.Endpoint)
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimSuffix(endpoint, "/")
	endpoint = strings.TrimPrefix(endpoint, strings.TrimSpace(d.Bucket)+".")
	return endpoint
}

func (d *ossDriver) bucketClient() (*oss.Bucket, error) {
	client, err := oss.New(d.normalizedEndpoint(), d.AccessKeyID, d.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("初始化阿里云OSS客户端失败: %w", err)
	}
	bucket, err := client.Bucket(d.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取阿里云OSS Bucket失败: %w", err)
	}
	return bucket, nil
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

type storageOssPlugin struct {
	meta   plugin.Meta
	driver *ossDriver
}

func init() {
	var m plugin.Meta
	json.Unmarshal(metaJSON, &m)
	plugin.Register(&storageOssPlugin{meta: m, driver: &ossDriver{}})
}

func (p *storageOssPlugin) Meta() plugin.Meta                    { return p.meta }
func (p *storageOssPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (p *storageOssPlugin) Migrate(_ *gorm.DB) error             { return nil }
func (p *storageOssPlugin) Install() error {
	kv := func(k string) string {
		var c model.ConfigKV
		if db.DB.Where("plugin=? AND key=?", "storage_oss", k).First(&c).Error == nil {
			return c.Value
		}
		return ""
	}
	p.driver.Endpoint = kv("endpoint")
	p.driver.Bucket = kv("bucket")
	p.driver.AccessKeyID = kv("access_key_id")
	p.driver.AccessKeySecret = kv("access_key_secret")
	driverStorage.Register(p.driver)
	return nil
}
func (p *storageOssPlugin) Uninstall() error { return nil }
