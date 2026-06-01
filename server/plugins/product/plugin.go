package product

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	productapi "github.com/ijry/lyshop/plugins/product/api"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
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
	if err := db.AutoMigrate(
		&productmodel.ProductCategory{},
		&productmodel.Product{},
		&productmodel.ProductSku{},
		&productmodel.ProductImage{},
		&productmodel.ProductFavorite{},
		&productmodel.SpecTemplate{},
	); err != nil {
		return err
	}
	migrator := db.Migrator()
	if migrator.HasIndex(&productmodel.ProductSku{}, "uk_product_sku_key") {
		_ = migrator.DropIndex(&productmodel.ProductSku{}, "uk_product_sku_key")
	}
	return migrator.CreateIndex(&productmodel.ProductSku{}, "uk_product_sku_key")
}

func (p *productPlugin) Install() error   { return nil }
func (p *productPlugin) Uninstall() error { return nil }
