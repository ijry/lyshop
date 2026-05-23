package message

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	msgapi "github.com/ijry/lyshop/plugins/message/api"
	msgmodel "github.com/ijry/lyshop/plugins/message/model"
	"github.com/ijry/lyshop/core/plugin"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type messagePlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("message plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&messagePlugin{meta: m})
}

func (p *messagePlugin) Meta() plugin.Meta { return p.meta }

func (p *messagePlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	msgapi.RegisterFrontRoutes(front)
	msgapi.RegisterAdminRoutes(admin)
}

func (p *messagePlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&msgmodel.Message{})
}

func (p *messagePlugin) Install() error   { return nil }
func (p *messagePlugin) Uninstall() error { return nil }
