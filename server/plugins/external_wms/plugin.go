package external_wms

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	externalapi "github.com/ijry/lyshop/plugins/external_wms/api"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type externalPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("external_wms plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&externalPlugin{meta: m})
}

func (p *externalPlugin) Meta() plugin.Meta { return p.meta }

func (p *externalPlugin) RegisterRoutes(_ *gin.RouterGroup, admin *gin.RouterGroup) {
	externalapi.RegisterAdminRoutes(admin)
}

func (p *externalPlugin) Migrate(_ *gorm.DB) error { return nil }
func (p *externalPlugin) Install() error           { return nil }
func (p *externalPlugin) Uninstall() error         { return nil }
