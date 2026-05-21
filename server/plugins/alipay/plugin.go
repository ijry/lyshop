package alipay

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/driver/payment"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type alipayPlugin struct {
	meta   plugin.Meta
	driver *AlipayDriver
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("alipay plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&alipayPlugin{meta: m, driver: &AlipayDriver{}})
}

func (p *alipayPlugin) Meta() plugin.Meta { return p.meta }

func (p *alipayPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	front.POST("/notify/alipay", func(c *gin.Context) {
		result, err := p.driver.HandleNotify(c.Request.Context(), c.Request)
		if err != nil || !result.Paid {
			c.String(500, "fail")
			return
		}
		// TODO: update order status to paid
		c.String(200, "success")
	})
}

func (p *alipayPlugin) Migrate(_ *gorm.DB) error { return nil }

func (p *alipayPlugin) Install() error {
	loadKV := func(key string) string {
		var cfg model.ConfigKV
		if err := db.DB.Where("plugin = ? AND key = ?", "alipay", key).First(&cfg).Error; err == nil {
			return cfg.Value
		}
		return ""
	}
	p.driver.AppID = loadKV("app_id")
	p.driver.PrivateKey = loadKV("private_key")
	p.driver.PublicKey = loadKV("public_key")
	p.driver.Sandbox = loadKV("sandbox") == "true"
	payment.Register(p.driver)
	return nil
}

func (p *alipayPlugin) Uninstall() error { return nil }
