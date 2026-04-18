package database

import (
	"fmt"
	"path/filepath"

	"halo/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Open(path string) (*gorm.DB, error) {
	dsn := filepath.ToSlash(filepath.Clean(path)) + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.UserRegistration{},
		&model.PasswordReset{},
		&model.PolicyRule{},
		&model.Role{},
		&model.RoleRule{},
		&model.UserRole{},
		&model.Post{},
		&model.Page{},
		&model.Category{},
		&model.Tag{},
		&model.Menu{},
		&model.Comment{},
		&model.Reply{},
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
