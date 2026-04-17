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
