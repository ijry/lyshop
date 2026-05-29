package wms

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	wmsapi "github.com/ijry/lyshop/plugins/wms/api"
	wmsmodel "github.com/ijry/lyshop/plugins/wms/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type wmsPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("wms plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&wmsPlugin{meta: m})
}

func (p *wmsPlugin) Meta() plugin.Meta { return p.meta }

func (p *wmsPlugin) RegisterRoutes(_ *gin.RouterGroup, admin *gin.RouterGroup) {
	wmsapi.RegisterAdminRoutes(admin)
}

func (p *wmsPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&wmsmodel.Warehouse{},
		&wmsmodel.InventoryStock{},
		&wmsmodel.InventoryMovement{},
		&wmsmodel.InventoryDoc{},
		&wmsmodel.InventoryDocItem{},
		&wmsmodel.InventoryReservation{},
	)
}

func (p *wmsPlugin) Install() error   { return nil }
func (p *wmsPlugin) Uninstall() error { return nil }
