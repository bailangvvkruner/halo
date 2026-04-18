package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type PageService struct {
	db *gorm.DB
}

func NewPageService(db *gorm.DB) *PageService {
	return &PageService{db: db}
}

func (s *PageService) List() ([]model.Page, error) {
	var items []model.Page
	return items, s.db.Order("id desc").Find(&items).Error
}

func (s *PageService) Create(page *model.Page) error {
	return s.db.Create(page).Error
}

func (s *PageService) Get(id uint) (*model.Page, error) {
	var page model.Page
	if err := s.db.First(&page, id).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

func (s *PageService) Update(id uint, page *model.Page) error {
	return s.db.Model(&model.Page{}).Where("id = ?", id).Updates(page).Error
}

func (s *PageService) Delete(id uint) error {
	return s.db.Delete(&model.Page{}, id).Error
}

func (s *PageService) Search(keyword string, result *[]model.Page) error {
	return s.db.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(result).Error
}
