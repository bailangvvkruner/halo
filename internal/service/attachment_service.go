package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type AttachmentService struct {
	db *gorm.DB
}

func NewAttachmentService(db *gorm.DB) *AttachmentService {
	return &AttachmentService{db: db}
}

func (s *AttachmentService) List() ([]model.Attachment, error) {
	var items []model.Attachment
	return items, s.db.Order("id desc").Find(&items).Error
}
