package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type SettingService struct {
	db *gorm.DB
}

func NewSettingService(db *gorm.DB) *SettingService {
	return &SettingService{db: db}
}

func (s *SettingService) List() ([]model.Setting, error) {
	var items []model.Setting
	return items, s.db.Order("id desc").Find(&items).Error
}
