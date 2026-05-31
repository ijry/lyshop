package im

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	imapi "github.com/ijry/lyshop/plugins/im/api"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	imsvc "github.com/ijry/lyshop/plugins/im/service"
	"github.com/ijry/lyshop/core/plugin"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type imPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("im plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&imPlugin{meta: m})
}

func (p *imPlugin) Meta() plugin.Meta { return p.meta }

func (p *imPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	imapi.RegisterFrontRoutes(front)
	imapi.RegisterAdminRoutes(admin)
	// WebSocket route is registered on the root engine via app.go
}

func (p *imPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&immodel.ImStaff{},
		&immodel.ImSession{},
		&immodel.ImMessage{},
		&immodel.ImAutoReply{},
		&immodel.ImTransferLog{},
		&immodel.ImKnowledge{},
	)
}

func (p *imPlugin) Install() error {
	go imsvc.GlobalHub.Run()
	return nil
}

func (p *imPlugin) Uninstall() error { return nil }
