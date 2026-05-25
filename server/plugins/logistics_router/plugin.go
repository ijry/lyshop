package logistics_router

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	logisticsRouterSvc "github.com/ijry/lyshop/plugins/logistics_router/service"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type logisticsRouterPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("logistics_router plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&logisticsRouterPlugin{meta: m})
}

func (p *logisticsRouterPlugin) Meta() plugin.Meta { return p.meta }

func (p *logisticsRouterPlugin) RegisterRoutes(_, _ *gin.RouterGroup) {}

func (p *logisticsRouterPlugin) Migrate(_ *gorm.DB) error { return nil }

func (p *logisticsRouterPlugin) Install() error {
	logisticsRouterSvc.ApplyRouteConfig()
	logisticsRouterSvc.StartPollingLoop()
	return nil
}

func (p *logisticsRouterPlugin) Uninstall() error { return nil }
