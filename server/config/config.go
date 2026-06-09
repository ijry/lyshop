package config

import "github.com/spf13/viper"

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Plugins     PluginsConfig     `mapstructure:"plugins"`
	Inventory   InventoryConfig   `mapstructure:"inventory"`
	ExternalWMS ExternalWMSConfig `mapstructure:"external_wms"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	DSN     string `mapstructure:"dsn"`
	MaxOpen int    `mapstructure:"max_open"`
	MaxIdle int    `mapstructure:"max_idle"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type PluginsConfig struct {
	Enabled []string `mapstructure:"enabled"`
}

type InventoryConfig struct {
	Provider     string `mapstructure:"provider"`
	ExternalMode string `mapstructure:"external_mode"`
}

type ExternalWMSConfig struct {
	Endpoint        string                 `mapstructure:"endpoint"`
	AppKey          string                 `mapstructure:"app_key"`
	AppSecret       string                 `mapstructure:"app_secret"`
	TimeoutMS       int                    `mapstructure:"timeout_ms"`
	CallbackEnabled bool                   `mapstructure:"callback_enabled"`
	SignatureTTL    int                    `mapstructure:"signature_ttl"`
	WorkerIntervalSec int                  `mapstructure:"worker_interval_sec"`
	Retry           ExternalWMSRetryConfig `mapstructure:"retry"`
}

type ExternalWMSRetryConfig struct {
	MaxAttempts    int `mapstructure:"max_attempts"`
	BackoffSeconds int `mapstructure:"backoff_seconds"`
}

// Global is the loaded config, available after Load().
var Global Config

// Load reads the YAML file at path and unmarshals it into Global.
func Load(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return viper.Unmarshal(&Global)
}
