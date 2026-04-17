package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type TagService struct {
	db *gorm.DB
}

func NewTagService(db *gorm.DB) *TagService {
	return &TagService{db: db}
}

func (s *TagService) List() ([]model.Tag, error) {
	var items []model.Tag
	return items, s.db.Order("id desc").Find(&items).Error
}
