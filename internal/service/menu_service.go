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

func (s *MenuService) Get(id uint) (*model.Menu, error) {
	var item model.Menu
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *MenuService) Create(menu *model.Menu) error {
	return s.db.Create(menu).Error
}

func (s *MenuService) Update(id uint, menu *model.Menu) error {
	return s.db.Model(&model.Menu{}).Where("id = ?", id).Updates(menu).Error
}

func (s *MenuService) Delete(id uint) error {
	return s.db.Delete(&model.Menu{}, id).Error
}
