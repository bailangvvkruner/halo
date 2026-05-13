package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.Exec("PRAGMA journal_mode=WAL").Error; err != nil {
		return nil, err
	}
	if err := db.Exec("PRAGMA foreign_keys=ON").Error; err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&extensionRecord{}); err != nil {
		return nil, err
	}

	return db, nil
}