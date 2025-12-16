package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Driver string // "sqlite" or "postgres"
	DSN    string // Data Source Name
}

// New creates a new database connection
func New(cfg Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// GetConfigFromEnv reads database config from environment variables
func GetConfigFromEnv() Config {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite" // default to sqlite
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" && driver == "sqlite" {
		dsn = "sticktoss.db" // default sqlite file
	}

	return Config{
		Driver: driver,
		DSN:    dsn,
	}
}
