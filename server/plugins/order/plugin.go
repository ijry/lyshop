package order

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	orderapi "github.com/ijry/lyshop/plugins/order/api"
	ordermodel "github.com/ijry/lyshop/plugins/order/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type orderPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("order plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&orderPlugin{meta: m})
}

func (p *orderPlugin) Meta() plugin.Meta { return p.meta }

func (p *orderPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	orderapi.RegisterFrontRoutes(front)
	orderapi.RegisterAdminRoutes(admin)
}

func (p *orderPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&ordermodel.Address{},
		&ordermodel.Order{},
		&ordermodel.OrderItem{},
		&ordermodel.OrderPayment{},
		&ordermodel.OrderRefund{},
		&ordermodel.OrderReview{},
		&ordermodel.OrderReviewAppend{},
		&ordermodel.OrderReviewReply{},
	)
}

func (p *orderPlugin) Install() error   { return nil }
func (p *orderPlugin) Uninstall() error { return nil }
