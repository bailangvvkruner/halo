package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

func (s *CategoryService) List() ([]model.Category, error) {
	var items []model.Category
	return items, s.db.Order("id desc").Find(&items).Error
}
