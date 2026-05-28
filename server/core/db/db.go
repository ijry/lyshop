package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global GORM instance, available after Init().
var DB *gorm.DB

const defaultSQLiteDSN = "lyshop.db"

func Init() error {
	cfg := config.Global.Database
	logLevel := logger.Silent
	if config.Global.Server.Mode == "debug" {
		logLevel = logger.Info
	}

	driverName, dsn := resolveDatabaseTarget(cfg.DSN)
	var dialector gorm.Dialector
	if driverName == "sqlite" {
		dialector = sqlite.Open(dsn)
	} else {
		dialector = mysql.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("db connect (%s): %w", driverName, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("db pool (%s): %w", driverName, err)
	}

	maxOpen := cfg.MaxOpen
	maxIdle := cfg.MaxIdle
	if driverName == "sqlite" {
		if maxOpen <= 0 {
			maxOpen = 1
		}
		if maxIdle <= 0 || maxIdle > maxOpen {
			maxIdle = maxOpen
		}
	}
	if maxOpen > 0 {
		sqlDB.SetMaxOpenConns(maxOpen)
	}
	if maxIdle > 0 {
		sqlDB.SetMaxIdleConns(maxIdle)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}

func resolveDatabaseTarget(rawDSN string) (driverName, dsn string) {
	trimmed := strings.TrimSpace(rawDSN)
	if trimmed == "" {
		return "sqlite", defaultSQLiteDSN
	}

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "sqlite://") {
		sqliteDSN := strings.TrimSpace(trimmed[len("sqlite://"):])
		if sqliteDSN == "" {
			return "sqlite", defaultSQLiteDSN
		}
		return "sqlite", sqliteDSN
	}
	if strings.HasPrefix(lower, "sqlite:") {
		sqliteDSN := strings.TrimSpace(trimmed[len("sqlite:"):])
		if sqliteDSN == "" {
			return "sqlite", defaultSQLiteDSN
		}
		return "sqlite", sqliteDSN
	}

	if strings.HasPrefix(lower, "file:") ||
		lower == ":memory:" ||
		strings.Contains(lower, "mode=memory") ||
		strings.HasSuffix(lower, ".db") ||
		strings.Contains(lower, ".db?") ||
		strings.HasSuffix(lower, ".sqlite") ||
		strings.Contains(lower, ".sqlite?") ||
		strings.HasSuffix(lower, ".sqlite3") ||
		strings.Contains(lower, ".sqlite3?") {
		return "sqlite", trimmed
	}

	return "mysql", trimmed
}
