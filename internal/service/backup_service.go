package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type BackupService struct {
	db *gorm.DB
}

func NewBackupService(db *gorm.DB) *BackupService {
	return &BackupService{db: db}
}

func (s *BackupService) List() ([]model.Backup, error) {
	var items []model.Backup
	return items, s.db.Order("id desc").Find(&items).Error
}
