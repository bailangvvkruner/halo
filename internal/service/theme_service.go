package service

import (
	"fmt"
	"os"
	"path/filepath"

	"halo/internal/config"
	"halo/internal/model"

	"gorm.io/gorm"
)

type ThemeService struct {
	db  *gorm.DB
	cfg config.Config
}

func NewThemeService(db *gorm.DB, cfg config.Config) *ThemeService {
	return &ThemeService{db: db, cfg: cfg}
}

func (s *ThemeService) List() ([]model.Theme, error) {
	var items []model.Theme
	return items, s.db.Order("id desc").Find(&items).Error
}

func (s *ThemeService) Activate(id uint) error {
	if err := s.db.Model(&model.Theme{}).Where("activated = ?", true).Update("activated", false).Error; err != nil {
		return err
	}
	return s.db.Model(&model.Theme{}).Where("id = ?", id).Update("activated", true).Error
}

func (s *ThemeService) ScanDirectories() ([]model.Theme, error) {
	themesDir := filepath.Join(s.cfg.WorkDir, "themes")
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, err
	}

	result := make([]model.Theme, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		result = append(result, model.Theme{
			Name:        entry.Name(),
			DisplayName: fmt.Sprintf("%s Theme", entry.Name()),
		})
	}

	return result, nil
}

func (s *ThemeService) SyncScannedThemes() ([]model.Theme, error) {
	items, err := s.ScanDirectories()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		var count int64
		s.db.Model(&model.Theme{}).Where("name = ?", item.Name).Count(&count)
		if count == 0 {
			copy := item
			if err := s.db.Create(&copy).Error; err != nil {
				return nil, err
			}
		}
	}

	return s.List()
}
