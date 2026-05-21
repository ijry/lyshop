package decor

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	decorapi "github.com/ijry/lyshop/plugins/decor/api"
	decormodel "github.com/ijry/lyshop/plugins/decor/model"
	"github.com/ijry/lyshop/core/plugin"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type decorPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("decor plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&decorPlugin{meta: m})
}

func (p *decorPlugin) Meta() plugin.Meta { return p.meta }

func (p *decorPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	decorapi.RegisterFrontRoutes(front)
	decorapi.RegisterAdminRoutes(admin)
}

func (p *decorPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&decormodel.DecorPage{})
}

func (p *decorPlugin) Install() error   { return nil }
func (p *decorPlugin) Uninstall() error { return nil }
