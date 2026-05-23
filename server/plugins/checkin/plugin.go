package checkin

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	checkinapi "github.com/ijry/lyshop/plugins/checkin/api"
	checkinmodel "github.com/ijry/lyshop/plugins/checkin/model"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/plugin"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type checkinPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("checkin plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&checkinPlugin{meta: m})
}

func (p *checkinPlugin) Meta() plugin.Meta { return p.meta }

func (p *checkinPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	checkinapi.RegisterFrontRoutes(front)
	checkinapi.RegisterAdminRoutes(admin)
}

func (p *checkinPlugin) Migrate(d *gorm.DB) error {
	return d.AutoMigrate(&checkinmodel.CheckinRule{}, &checkinmodel.CheckinLog{})
}

func (p *checkinPlugin) Install() error {
	// Seed default rules if empty
	var count int64
	db.DB.Model(&checkinmodel.CheckinRule{}).Count(&count)
	if count == 0 {
		db.DB.Create(&[]checkinmodel.CheckinRule{
			{Day: 0, Points: 10},  // default: 10 points
			{Day: 3, Points: 20},  // 3rd consecutive day: 20 points
			{Day: 7, Points: 50},  // 7th consecutive day: 50 points
		})
	}
	return nil
}

func (p *checkinPlugin) Uninstall() error { return nil }
