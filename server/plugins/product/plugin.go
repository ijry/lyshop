package product

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	productapi "github.com/ijry/lyshop/plugins/product/api"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
	"github.com/ijry/lyshop/core/plugin"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type productPlugin struct {
	meta plugin.Meta
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("product plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&productPlugin{meta: m})
}

func (p *productPlugin) Meta() plugin.Meta { return p.meta }

func (p *productPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	productapi.RegisterFrontRoutes(front)
	productapi.RegisterAdminRoutes(admin)
}

func (p *productPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&productmodel.ProductCategory{},
		&productmodel.Product{},
		&productmodel.ProductSku{},
		&productmodel.ProductImage{},
	)
}

func (p *productPlugin) Install() error   { return nil }
func (p *productPlugin) Uninstall() error { return nil }
