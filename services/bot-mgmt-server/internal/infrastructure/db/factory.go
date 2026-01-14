package db

import (
	"fmt"

	"github.com/SKD-fastcampus/bot-management/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB creates a new database connection based on configuration
func NewDB(cfg config.Config) (*gorm.DB, error) {
	driver := cfg.GetString("db.driver")
	dsn := ""

	var dialector gorm.Dialector

	switch driver {
	case "mysql":
		// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.GetString("db.user"),
			cfg.GetString("db.password"),
			cfg.GetString("db.host"),
			cfg.GetString("db.port"),
			cfg.GetString("db.name"),
		)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
			cfg.GetString("db.host"),
			cfg.GetString("db.user"),
			cfg.GetString("db.password"),
			cfg.GetString("db.name"),
			cfg.GetString("db.port"),
		)
		dialector = postgres.Open(dsn)
	case "":
		return nil, fmt.Errorf("database driver is not specified in config (db.driver)")
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database using %s driver: %w", driver, err)
	}

	return db, nil
}
