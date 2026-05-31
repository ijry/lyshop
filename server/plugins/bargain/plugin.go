package bargain

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/plugins/bargain/api"
	"github.com/ijry/lyshop/plugins/bargain/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var pluginJSON []byte

type Plugin struct{}

func (p *Plugin) Meta() plugin.Meta {
	var meta plugin.Meta
	_ = json.Unmarshal(pluginJSON, &meta)
	return meta
}

func (p *Plugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	api.RegisterAdminRoutes(admin)
	api.RegisterFrontRoutes(front)
}

func (p *Plugin) Migrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&model.BargainActivity{},
		&model.BargainProduct{},
		&model.BargainOrder{},
		&model.BargainHelper{},
	)
}

func (p *Plugin) Install() error {
	return p.Migrate(db.DB)
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
