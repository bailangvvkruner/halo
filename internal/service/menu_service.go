package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type MenuService struct {
	db *gorm.DB
}

func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{db: db}
}

func (s *MenuService) List() ([]model.Menu, error) {
	var items []model.Menu
	return items, s.db.Order("id desc").Find(&items).Error
}
