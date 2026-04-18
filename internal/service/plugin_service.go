package service

import (
	"fmt"
	"os"
	"path/filepath"

	"halo/internal/config"
	"halo/internal/model"

	"gorm.io/gorm"
)

type PluginService struct {
	db  *gorm.DB
	cfg config.Config
}

func NewPluginService(db *gorm.DB, cfg config.Config) *PluginService {
	return &PluginService{db: db, cfg: cfg}
}

func (s *PluginService) List() ([]model.Plugin, error) {
	var items []model.Plugin
	return items, s.db.Order("id desc").Find(&items).Error
}

func (s *PluginService) Get(id uint) (*model.Plugin, error) {
	var item model.Plugin
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *PluginService) Toggle(id uint, enabled bool) error {
	return s.db.Model(&model.Plugin{}).Where("id = ?", id).Update("enabled", enabled).Error
}

func (s *PluginService) ScanDirectories() ([]model.Plugin, error) {
	pluginsDir := filepath.Join(s.cfg.WorkDir, "plugins")
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return nil, err
	}

	result := make([]model.Plugin, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		result = append(result, model.Plugin{
			Name:        entry.Name(),
			DisplayName: fmt.Sprintf("%s Plugin", entry.Name()),
			Path:        filepath.Join(pluginsDir, entry.Name()),
		})
	}

	return result, nil
}
