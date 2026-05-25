package delivery_express

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	deliveryDriver "github.com/ijry/lyshop/core/driver/delivery"
	"github.com/ijry/lyshop/core/plugin"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type deliveryExpressPlugin struct {
	meta   plugin.Meta
	driver *expressDriver
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("delivery_express plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&deliveryExpressPlugin{meta: m, driver: &expressDriver{}})
}

func (p *deliveryExpressPlugin) Meta() plugin.Meta { return p.meta }

func (p *deliveryExpressPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}

func (p *deliveryExpressPlugin) Migrate(_ *gorm.DB) error { return nil }

func (p *deliveryExpressPlugin) Install() error {
	deliveryDriver.Register(p.driver)
	return nil
}

func (p *deliveryExpressPlugin) Uninstall() error { return nil }
