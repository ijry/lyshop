package logistics_kdniao

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

type logisticsKdniaoPlugin struct {
	meta   plugin.Meta
	driver *kdniaoDriver
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("logistics_kdniao plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&logisticsKdniaoPlugin{meta: m, driver: &kdniaoDriver{}})
}

func (p *logisticsKdniaoPlugin) Meta() plugin.Meta                    { return p.meta }
func (p *logisticsKdniaoPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}
func (p *logisticsKdniaoPlugin) Migrate(_ *gorm.DB) error             { return nil }

func (p *logisticsKdniaoPlugin) Install() error {
	loadKV := func(key, defaultValue string) string {
		var cfg model.ConfigKV
		if err := db.DB.Where("plugin = ? AND key = ?", "logistics_kdniao", key).First(&cfg).Error; err == nil {
			if cfg.Value != "" {
				return cfg.Value
			}
		}
		return defaultValue
	}
	p.driver.APIURL = loadKV("api_url", "https://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx")
	p.driver.EBusinessID = loadKV("ebusiness_id", "")
	p.driver.APIKey = loadKV("api_key", "")
	logisticsDriver.Register(p.driver)
	return nil
}

func (p *logisticsKdniaoPlugin) Uninstall() error { return nil }
