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

func (s *SettingService) Get(id uint) (*model.Setting, error) {
	var item model.Setting
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *SettingService) Update(id uint, setting *model.Setting) error {
	return s.db.Model(&model.Setting{}).Where("id = ?", id).Updates(setting).Error
}
