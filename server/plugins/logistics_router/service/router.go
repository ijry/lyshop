package service

import (
	"strconv"
	"strings"

	"github.com/ijry/lyshop/core/db"
	logisticsDriver "github.com/ijry/lyshop/core/driver/logistics"
	"github.com/ijry/lyshop/model"
)

func loadCfg(key, defaultValue string) string {
	var cfg model.ConfigKV
	if err := db.DB.Where("plugin = ? AND key = ?", "logistics_router", key).First(&cfg).Error; err == nil {
		if value := strings.TrimSpace(cfg.Value); value != "" {
			return value
		}
	}
	return defaultValue
}

func loadBool(key string, defaultValue bool) bool {
	value := strings.ToLower(strings.TrimSpace(loadCfg(key, "")))
	if value == "" {
		return defaultValue
	}
	return value == "1" || value == "true" || value == "on" || value == "yes"
}

func loadInt(key string, defaultValue int) int {
	value := strings.TrimSpace(loadCfg(key, ""))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func ApplyRouteConfig() {
	primary := loadCfg("primary_driver", "kuaidi100")
	secondary := loadCfg("secondary_driver", "kdniao")
	logisticsDriver.SetDefaultDrivers(primary, secondary)
}
