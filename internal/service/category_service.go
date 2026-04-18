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

func (s *CategoryService) Get(id uint) (*model.Category, error) {
	var item model.Category
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *CategoryService) Create(category *model.Category) error {
	return s.db.Create(category).Error
}

func (s *CategoryService) Update(id uint, category *model.Category) error {
	return s.db.Model(&model.Category{}).Where("id = ?", id).Updates(category).Error
}

func (s *CategoryService) Delete(id uint) error {
	return s.db.Delete(&model.Category{}, id).Error
}
