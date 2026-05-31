package group_buy

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/plugins/group_buy/api"
	"github.com/ijry/lyshop/plugins/group_buy/model"
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
		&model.GroupBuyActivity{},
		&model.GroupBuyProduct{},
		&model.GroupBuyOrder{},
		&model.GroupBuyMember{},
	)
}

func (p *Plugin) Install() error {
	return p.Migrate(db.DB)
}

func (p *Plugin) Uninstall() error {
	return db.DB.Migrator().DropTable(
		&model.GroupBuyActivity{},
		&model.GroupBuyProduct{},
		&model.GroupBuyOrder{},
		&model.GroupBuyMember{},
	)
}

func init() {
	plugin.Register(&Plugin{})
}
