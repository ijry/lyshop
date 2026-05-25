package vip

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	vipapi "github.com/ijry/lyshop/plugins/vip/api"
	_ "github.com/ijry/lyshop/plugins/vip/calculator"
	vipmodel "github.com/ijry/lyshop/plugins/vip/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type vipPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("vip plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&vipPlugin{meta: m})
}

func (p *vipPlugin) Meta() plugin.Meta { return p.meta }

func (p *vipPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	vipapi.RegisterFrontRoutes(front)
	vipapi.RegisterAdminRoutes(admin)
}

func (p *vipPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&vipmodel.Plan{},
		&vipmodel.Level{},
		&vipmodel.UserAsset{},
		&vipmodel.GrowthLog{},
		&vipmodel.CouponRule{},
		&vipmodel.CouponClaim{},
		&vipmodel.SkuPrice{},
		&vipmodel.OrderBenefit{},
	)
}

func (p *vipPlugin) Install() error   { return nil }
func (p *vipPlugin) Uninstall() error { return nil }
