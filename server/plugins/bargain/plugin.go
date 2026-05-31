package bargain

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/server/plugins/bargain/api"
	"github.com/ijry/lyshop/server/plugins/bargain/model"
)

//go:embed plugin.json
var pluginJSON []byte

type Plugin struct{}

func (p *Plugin) Meta() plugin.Metadata {
	var meta plugin.Metadata
	_ = json.Unmarshal(pluginJSON, &meta)
	return meta
}

func (p *Plugin) RegisterRoutes(r *gin.Engine) {
	api.RegisterAdminRoutes(r.Group(""))
	api.RegisterFrontRoutes(r.Group(""))
}

func (p *Plugin) Migrate() error {
	return db.DB.AutoMigrate(
		&model.BargainActivity{},
		&model.BargainProduct{},
		&model.BargainOrder{},
		&model.BargainHelper{},
	)
}

func (p *Plugin) Install() error {
	return p.Migrate()
}

func (p *Plugin) Uninstall() error {
	return db.DB.Migrator().DropTable(
		&model.BargainActivity{},
		&model.BargainProduct{},
		&model.BargainOrder{},
		&model.BargainHelper{},
	)
}

func init() {
	plugin.Register(&Plugin{})
}
