package database

import (
	"fmt"
	"os"
	"path/filepath"

	"halo/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(path string) (*gorm.DB, error) {
	if _, err := os.Stat(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("stat database: %w", err)
	}

	if os.Getenv("CGO_ENABLED") == "0" {
		return nil, fmt.Errorf("sqlite runtime requires cgo-enabled Go environment on this machine")
	}

	dsn := filepath.ToSlash(filepath.Clean(path)) + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Page{},
		&model.Category{},
		&model.Tag{},
		&model.Menu{},
		&model.Comment{},
		&model.Theme{},
		&model.Plugin{},
		&model.Attachment{},
		&model.Backup{},
		&model.Setting{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return db, nil
}
