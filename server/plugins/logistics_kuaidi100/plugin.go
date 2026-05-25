package logistics_kuaidi100

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	logisticsDriver "github.com/ijry/lyshop/core/driver/logistics"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type logisticsKuaidi100Plugin struct {
	meta   plugin.Meta
	driver *kuaidi100Driver
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("logistics_kuaidi100 plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&logisticsKuaidi100Plugin{meta: m, driver: &kuaidi100Driver{}})
}

func (p *logisticsKuaidi100Plugin) Meta() plugin.Meta                    { return p.meta }
func (p *logisticsKuaidi100Plugin) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (p *logisticsKuaidi100Plugin) Migrate(_ *gorm.DB) error             { return nil }

func (p *logisticsKuaidi100Plugin) Install() error {
	loadKV := func(key, defaultValue string) string {
		var cfg model.ConfigKV
		if err := db.DB.Where("plugin = ? AND key = ?", "logistics_kuaidi100", key).First(&cfg).Error; err == nil {
			if cfg.Value != "" {
				return cfg.Value
			}
		}
		return defaultValue
	}
	p.driver.APIURL = loadKV("api_url", "https://poll.kuaidi100.com/poll/query.do")
	p.driver.Customer = loadKV("customer", "")
	p.driver.Key = loadKV("key", "")
	logisticsDriver.Register(p.driver)
	return nil
}

func (p *logisticsKuaidi100Plugin) Uninstall() error { return nil }
