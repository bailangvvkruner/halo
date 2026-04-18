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

func (s *TagService) Get(id uint) (*model.Tag, error) {
	var item model.Tag
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *TagService) Create(tag *model.Tag) error {
	return s.db.Create(tag).Error
}

func (s *TagService) Update(id uint, tag *model.Tag) error {
	return s.db.Model(&model.Tag{}).Where("id = ?", id).Updates(tag).Error
}

func (s *TagService) Delete(id uint) error {
	return s.db.Delete(&model.Tag{}, id).Error
}
