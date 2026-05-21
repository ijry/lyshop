package sms

import (
	_ "embed"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	driverSms "github.com/ijry/lyshop/core/driver/sms"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/core/response"
	"github.com/ijry/lyshop/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type smsPlugin struct {
	meta   plugin.Meta
	driver *SmsDriver
}

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("sms plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&smsPlugin{meta: m, driver: &SmsDriver{}})
}

func (p *smsPlugin) Meta() plugin.Meta { return p.meta }

func (p *smsPlugin) RegisterRoutes(_ *gin.RouterGroup, admin *gin.RouterGroup) {
	admin.GET("/system/sms/config", func(c *gin.Context) {
		cfg := map[string]string{
			"provider":   p.driver.Provider,
			"access_key": p.driver.AccessKey,
			"sign_name":  p.driver.SignName,
		}
		response.OK(c, cfg)
	})
	admin.PUT("/system/sms/config", func(c *gin.Context) {
		var req struct {
			Provider  string `json:"provider"`
			AccessKey string `json:"access_key"`
			SecretKey string `json:"secret_key"`
			SignName  string `json:"sign_name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, 400, err.Error())
			return
		}
		saveKV := func(key, val string) {
			db.DB.Where(model.ConfigKV{Plugin: "sms", Key: key}).
				Assign(model.ConfigKV{Value: val}).
				FirstOrCreate(&model.ConfigKV{})
		}
		saveKV("provider", req.Provider)
		saveKV("access_key", req.AccessKey)
		saveKV("secret_key", req.SecretKey)
		saveKV("sign_name", req.SignName)

		p.driver.Provider = req.Provider
		p.driver.AccessKey = req.AccessKey
		p.driver.SecretKey = req.SecretKey
		p.driver.SignName = req.SignName
		response.OK(c, nil)
	})
}

func (p *smsPlugin) Migrate(_ *gorm.DB) error { return nil }

func (p *smsPlugin) Install() error {
	loadKV := func(key string) string {
		var cfg model.ConfigKV
		if err := db.DB.Where("plugin = ? AND key = ?", "sms", key).First(&cfg).Error; err == nil {
			return cfg.Value
		}
		return ""
	}
	p.driver.Provider = loadKV("provider")
	if p.driver.Provider == "" {
		p.driver.Provider = ProviderAliyun // default
	}
	p.driver.AccessKey = loadKV("access_key")
	p.driver.SecretKey = loadKV("secret_key")
	p.driver.SignName = loadKV("sign_name")
	driverSms.Register(p.driver)
	return nil
}

func (p *smsPlugin) Uninstall() error { return nil }
