package wechat_pay

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

type wechatPayPlugin struct {
	meta   plugin.Meta
	driver *WechatPayDriver
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("wechat_pay plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&wechatPayPlugin{meta: m, driver: &WechatPayDriver{}})
}

func (p *wechatPayPlugin) Meta() plugin.Meta { return p.meta }
func (p *wechatPayPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	// Notify callback — no auth
	front.POST("/notify/wechat", func(c *gin.Context) {
		result, err := p.driver.HandleNotify(c.Request.Context(), c.Request)
		if err != nil || !result.Paid {
			c.String(500, "FAIL")
			return
		}
		// TODO: update order status to paid
		c.String(200, "SUCCESS")
	})
}

func (p *wechatPayPlugin) Migrate(_ *gorm.DB) error { return nil }

func (p *wechatPayPlugin) Install() error {
	// Load config from DB
	var cfg model.ConfigKV
	if err := db.DB.Where("plugin = ? AND key = ?", "wechat_pay", "app_id").First(&cfg).Error; err == nil {
		p.driver.AppID = cfg.Value
	}
	var mch model.ConfigKV
	if err := db.DB.Where("plugin = ? AND key = ?", "wechat_pay", "mch_id").First(&mch).Error; err == nil {
		p.driver.MchID = mch.Value
	}
	var key model.ConfigKV
	if err := db.DB.Where("plugin = ? AND key = ?", "wechat_pay", "api_key").First(&key).Error; err == nil {
		p.driver.APIKey = key.Value
	}
	payment.Register(p.driver)
	return nil
}

func (p *wechatPayPlugin) Uninstall() error { return nil }
