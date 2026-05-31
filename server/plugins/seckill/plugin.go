package seckill

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/server/core/plugin"
	seckillapi "github.com/ijry/lyshop/server/plugins/seckill/api"
	_ "github.com/ijry/lyshop/server/plugins/seckill/calculator"
	seckillmodel "github.com/ijry/lyshop/server/plugins/seckill/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type seckillPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("seckill plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&seckillPlugin{meta: m})
}

func (p *seckillPlugin) Meta() plugin.Meta { return p.meta }

func (p *seckillPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	seckillapi.RegisterFrontRoutes(front)
	seckillapi.RegisterAdminRoutes(admin)
}

func (p *seckillPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&seckillmodel.SeckillActivity{},
		&seckillmodel.SeckillProduct{},
	)
}

func (p *seckillPlugin) Install() error   { return nil }
func (p *seckillPlugin) Uninstall() error { return nil }
