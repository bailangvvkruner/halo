package service

import (
	"archive/zip"
	"fmt"
	"halo/internal/model"
	"os"
	"path/filepath"
	"time"

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

func (s *BackupService) Get(id uint) (*model.Backup, error) {
	var item model.Backup
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *BackupService) Create(workDir string) (*model.Backup, error) {
	backupDir := filepath.Join(workDir, "backups")
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("halo-backup-%s.zip", time.Now().Format("20060102-150405"))
	fullPath := filepath.Join(backupDir, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	if _, err := zipWriter.Create("README.txt"); err != nil {
		zipWriter.Close()
		return nil, err
	}
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	stat, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	backup := &model.Backup{
		Filename: filename,
		Status:   "Succeeded",
		Size:     stat.Size(),
	}

	if err := s.db.Create(backup).Error; err != nil {
		return nil, err
	}

	return backup, nil
}
