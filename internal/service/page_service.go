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
