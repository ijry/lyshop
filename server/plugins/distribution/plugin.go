package distribution

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/plugins/distribution/api"
	"github.com/ijry/lyshop/plugins/distribution/model"
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
		&model.DistributionConfig{},
		&model.Distributor{},
		&model.DistributionOrder{},
		&model.DistributionWithdrawal{},
	)
}

func (p *Plugin) Install() error {
	return p.Migrate(db.DB)
}

func (p *Plugin) Uninstall() error {
	return db.DB.Migrator().DropTable(
		&model.DistributionConfig{},
		&model.Distributor{},
		&model.DistributionOrder{},
		&model.DistributionWithdrawal{},
	)
}

func init() {
	plugin.Register(&Plugin{})
}
