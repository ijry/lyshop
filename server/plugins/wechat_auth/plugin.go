package wechat_auth

import (
	_ "embed"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	driverOauth "github.com/ijry/lyshop/core/driver/oauth"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/core/response"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type wechatAuthPlugin struct {
	meta   plugin.Meta
	driver *WechatAuthDriver
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("wechat_auth plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&wechatAuthPlugin{meta: m, driver: &WechatAuthDriver{}})
}

func (p *wechatAuthPlugin) Meta() plugin.Meta { return p.meta }

func (p *wechatAuthPlugin) RegisterRoutes(front, _ *gin.RouterGroup) {
	// Mini-program login: POST /api/v1/auth/wechat/miniapp  body: {code}
	front.POST("/auth/wechat/miniapp", func(c *gin.Context) {
		var req struct {
			Code string `json:"code" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, 400, err.Error())
			return
		}
		info, err := p.driver.HandleCallback(c.Request.Context(), req.Code)
		if err != nil {
			response.Fail(c, 401, err.Error())
			return
		}
		// Find or create user by openid
		var user model.User
		err = db.DB.Where("phone = ?", "wx_"+info.OpenID).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = model.User{
				Phone:    "wx_" + info.OpenID,
				Nickname: info.Nickname,
				Avatar:   info.Avatar,
				Status:   1,
			}
			db.DB.Create(&user)
		}
		token, _ := middleware.GenerateToken(user.ID, "user", nil)
		response.OK(c, gin.H{"token": token, "open_id": info.OpenID})
	})
}

func (p *wechatAuthPlugin) Migrate(_ *gorm.DB) error { return nil }

func (p *wechatAuthPlugin) Install() error {
	loadKV := func(key string) string {
		var cfg model.ConfigKV
		if err := db.DB.Where("plugin = ? AND key = ?", "wechat_auth", key).First(&cfg).Error; err == nil {
			return cfg.Value
		}
		return ""
	}
	p.driver.MiniAppID = loadKV("mini_app_id")
	p.driver.MiniAppSecret = loadKV("mini_app_secret")
	p.driver.H5AppID = loadKV("h5_app_id")
	p.driver.H5AppSecret = loadKV("h5_app_secret")
	driverOauth.Register(p.driver)
	return nil
}

func (p *wechatAuthPlugin) Uninstall() error { return nil }
