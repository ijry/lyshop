package ai_image

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	aiapi "github.com/ijry/lyshop/plugins/ai_image/api"
	aidrv "github.com/ijry/lyshop/plugins/ai_image/driver"
	aimodel "github.com/ijry/lyshop/plugins/ai_image/model"
	"github.com/ijry/lyshop/core/db"
	aidriver "github.com/ijry/lyshop/core/driver/ai"
	"github.com/ijry/lyshop/core/plugin"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type aiImagePlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("ai_image plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&aiImagePlugin{meta: m})
}

func (p *aiImagePlugin) Meta() plugin.Meta { return p.meta }

func (p *aiImagePlugin) RegisterRoutes(_ *gin.RouterGroup, admin *gin.RouterGroup) {
	aiapi.RegisterAdminRoutes(admin)
}

func (p *aiImagePlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&aimodel.AiModel{}, &aimodel.AiImageTask{})
}

func (p *aiImagePlugin) Install() error {
	// Register all built-in AI drivers
	aidriver.Register(&aidrv.TongyiDriver{}, false)
	aidriver.Register(&aidrv.WenxinDriver{}, false)
	aidriver.Register(&aidrv.HunyuanDriver{}, false)
	aidriver.Register(&aidrv.OpenAIDriver{}, false)

	// Load API keys from DB and update drivers
	var models []aimodel.AiModel
	db.DB.Where("status = 1").Find(&models)
	for _, m := range models {
		isDefault := m.IsDefault == 1
		switch m.Driver {
		case "tongyi":
			aidriver.Register(&aidrv.TongyiDriver{APIKey: m.ApiKey}, isDefault)
		case "wenxin":
			aidriver.Register(&aidrv.WenxinDriver{APIKey: m.ApiKey}, isDefault)
		case "openai":
			aidriver.Register(&aidrv.OpenAIDriver{APIKey: m.ApiKey, Endpoint: m.Endpoint}, isDefault)
		}
	}
	return nil
}

func (p *aiImagePlugin) Uninstall() error { return nil }
