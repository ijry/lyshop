package delivery_router

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/core/response"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type deliveryRouterPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("delivery_router plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&deliveryRouterPlugin{meta: m})
}

func (p *deliveryRouterPlugin) Meta() plugin.Meta { return p.meta }

func (p *deliveryRouterPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	handler := func(c *gin.Context) {
		response.OK(c, gin.H{"mode": getDeliveryMode()})
	}
	admin.GET("/delivery/mode", handler)
	front.GET("/delivery/mode", handler)
}

func (p *deliveryRouterPlugin) Migrate(_ *gorm.DB) error { return nil }
func (p *deliveryRouterPlugin) Install() error            { return nil }
func (p *deliveryRouterPlugin) Uninstall() error          { return nil }

// getDeliveryMode reads the delivery_mode from ConfigKV, defaults to "express".
func getDeliveryMode() string {
	var kv model.ConfigKV
	if err := db.DB.Where("plugin = ? AND `key` = ?", "delivery_router", "delivery_mode").First(&kv).Error; err != nil {
		return "express"
	}
	mode := kv.Value
	if mode == "local" || mode == "both" {
		return mode
	}
	return "express"
}
