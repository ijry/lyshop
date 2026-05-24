package decor

import (
	_ "embed"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/core/db"
	"github.com/ijry/lyshop/core/plugin"
	decorapi "github.com/ijry/lyshop/plugins/decor/api"
	decormodel "github.com/ijry/lyshop/plugins/decor/model"
	"gorm.io/gorm"
)

//go:embed plugin.json
var metaJSON []byte

type decorPlugin struct{ meta plugin.Meta }

func init() {
	var m plugin.Meta
	if err := json.Unmarshal(metaJSON, &m); err != nil {
		panic("decor plugin: invalid plugin.json: " + err.Error())
	}
	plugin.Register(&decorPlugin{meta: m})
}

func (p *decorPlugin) Meta() plugin.Meta { return p.meta }

func (p *decorPlugin) RegisterRoutes(front, admin *gin.RouterGroup) {
	decorapi.RegisterFrontRoutes(front)
	decorapi.RegisterAdminRoutes(admin)
}

func (p *decorPlugin) Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&decormodel.DecorPage{}); err != nil {
		return err
	}
	m := db.Migrator()
	if m.HasIndex(&decormodel.DecorPage{}, "uk_merchant_page") {
		_ = m.DropIndex(&decormodel.DecorPage{}, "uk_merchant_page")
	}
	_ = db.Exec("UPDATE decor_pages SET variant_key = 'default' WHERE variant_key = '' OR variant_key IS NULL").Error
	_ = db.Exec("UPDATE decor_pages SET variant_name = '默认副本' WHERE variant_name = '' OR variant_name IS NULL").Error
	_ = db.Exec("UPDATE decor_pages SET is_published = CASE WHEN published_at IS NULL THEN 0 ELSE 1 END WHERE is_published IS NULL").Error
	return nil
}

func (p *decorPlugin) Install() error {
	type pageGroup struct {
		MerchantID uint64 `gorm:"column:merchant_id"`
		PageKey    string `gorm:"column:page_key"`
	}
	var groups []pageGroup
	if err := db.DB.Model(&decormodel.DecorPage{}).
		Select("merchant_id, page_key").
		Group("merchant_id, page_key").
		Find(&groups).Error; err != nil {
		return err
	}
	for _, g := range groups {
		_ = db.DB.Model(&decormodel.DecorPage{}).
			Where("merchant_id = ? AND page_key = ?", g.MerchantID, g.PageKey).
			Update("is_published", false).Error
		var latest decormodel.DecorPage
		err := db.DB.Where("merchant_id = ? AND page_key = ? AND published_at IS NOT NULL", g.MerchantID, g.PageKey).
			Order("published_at desc, id desc").
			First(&latest).Error
		if err == nil {
			_ = db.DB.Model(&decormodel.DecorPage{}).Where("id = ?", latest.ID).Update("is_published", true).Error
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	var pages []decormodel.DecorPage
	if err := db.DB.Where("variant_name = '' OR variant_name IS NULL").Find(&pages).Error; err == nil {
		for _, row := range pages {
			name := "默认副本"
			if key := strings.TrimSpace(row.VariantKey); key != "" && key != "default" {
				name = "副本 " + key
			}
			_ = db.DB.Model(&decormodel.DecorPage{}).Where("id = ?", row.ID).Update("variant_name", name).Error
		}
	}
	return nil
}
func (p *decorPlugin) Uninstall() error { return nil }
