package storage_router

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type storageRouterPlugin struct {
	meta plugin.Meta
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("storage_router plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&storageRouterPlugin{meta: m})
}

func (p *storageRouterPlugin) Meta() plugin.Meta { return p.meta }

func (p *storageRouterPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}

func (p *storageRouterPlugin) Migrate(_ *gorm.DB) error { return nil }

func (p *storageRouterPlugin) Install() error {
	var cfg model.ConfigKV
	if err := db.DB.Where("plugin = ? AND key = ?", "storage_router", "default_driver").First(&cfg).Error; err == nil {
		if selected := normalizeDriver(cfg.Value); selected != "" {
			driverStorage.SetDefault(selected)
			return nil
		}
	}
	if current := normalizeDriver(driverStorage.GetDefaultName()); current != "" {
		driverStorage.SetDefault(current)
		return nil
	}
	driverStorage.SetDefault("local")
	return nil
}

func (p *storageRouterPlugin) Uninstall() error { return nil }

func normalizeDriver(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "local", "storage_local":
		return "local"
	case "oss", "storage_oss", "aliyun_oss":
		return "oss"
	case "cos", "storage_cos", "qcloud_cos":
		return "cos"
	case "qiniu", "storage_qiniu", "qiniu_kodo":
		return "qiniu"
	default:
		return ""
	}
}
