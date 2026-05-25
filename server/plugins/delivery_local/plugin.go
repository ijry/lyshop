package delivery_local

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

type deliveryLocalPlugin struct {
	meta   plugin.Meta
	driver *localDriver
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("delivery_local plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&deliveryLocalPlugin{meta: m, driver: &localDriver{}})
}

func (p *deliveryLocalPlugin) Meta() plugin.Meta { return p.meta }

func (p *deliveryLocalPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}

func (p *deliveryLocalPlugin) Migrate(_ *gorm.DB) error { return nil }

func (p *deliveryLocalPlugin) Install() error {
	deliveryDriver.Register(p.driver)
	return nil
}

func (p *deliveryLocalPlugin) Uninstall() error { return nil }
