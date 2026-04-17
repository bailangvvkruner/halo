package database

import (
	"fmt"
	"path/filepath"

	"halo/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filepath.Clean(path)), &gorm.Config{})
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
