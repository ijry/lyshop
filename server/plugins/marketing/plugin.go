package marketing

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	mktapi "github.com/ijry/lyshop/plugins/marketing/api"
	mktmodel "github.com/ijry/lyshop/plugins/marketing/model"
	"github.com/ijry/lyshop/core/plugin"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type marketingPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("marketing plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&marketingPlugin{meta: m})
}

func (p *marketingPlugin) Meta() plugin.Meta { return p.meta }

func (p *marketingPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	mktapi.RegisterFrontRoutes(front)
	mktapi.RegisterAdminRoutes(admin)
}

func (p *marketingPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&mktmodel.Coupon{},
		&mktmodel.CouponUser{},
		&mktmodel.Activity{},
		&mktmodel.ActivityProduct{},
		&mktmodel.PointsLog{},
	)
}

func (p *marketingPlugin) Install() error   { return nil }
func (p *marketingPlugin) Uninstall() error { return nil }
