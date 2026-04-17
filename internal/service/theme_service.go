package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type ThemeService struct {
	db *gorm.DB
}

func NewThemeService(db *gorm.DB) *ThemeService {
	return &ThemeService{db: db}
}

func (s *ThemeService) List() ([]model.Theme, error) {
	var items []model.Theme
	return items, s.db.Order("id desc").Find(&items).Error
}
