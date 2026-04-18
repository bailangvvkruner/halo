package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type PluginService struct {
	db *gorm.DB
}

func NewPluginService(db *gorm.DB) *PluginService {
	return &PluginService{db: db}
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
