package points_mall

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/plugin"
	pmapi "github.com/ijry/lyshop/plugins/points_mall/api"
	_ "github.com/ijry/lyshop/plugins/points_mall/calculator"
	pmmodel "github.com/ijry/lyshop/plugins/points_mall/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type pointsMallPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("points_mall plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&pointsMallPlugin{meta: m})
}

func (p *pointsMallPlugin) Meta() plugin.Meta { return p.meta }

func (p *pointsMallPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	pmapi.RegisterFrontRoutes(front)
	pmapi.RegisterAdminRoutes(admin)
}

func (p *pointsMallPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&pmmodel.PointsLog{},
		&pmmodel.PointsProduct{},
		&pmmodel.PointsExchange{},
		&pmmodel.PointsConfig{},
	)
}

func (p *pointsMallPlugin) Install() error {
	// 初始化默认配置
	// TODO: 实现默认配置初始化
	return nil
}

func (p *pointsMallPlugin) Uninstall() error {
	return nil
}
