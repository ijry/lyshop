package db

import (
	"fmt"
	"time"

	"github.com/ijry/lyshop/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global GORM instance, available after Init().
var DB *gorm.DB

func Init() error {
	cfg := config.Global.Database
	logLevel := logger.Silent
	if config.Global.Server.Mode == "debug" {
		logLevel = logger.Info
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}
