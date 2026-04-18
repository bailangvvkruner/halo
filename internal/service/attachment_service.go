package service

import (
	"path/filepath"
	"strings"

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

func (s *AttachmentService) Create(attachment *model.Attachment) error {
	return s.db.Create(attachment).Error
}

func (s *AttachmentService) Get(id uint) (*model.Attachment, error) {
	var item model.Attachment
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *AttachmentService) Delete(id uint) error {
	return s.db.Delete(&model.Attachment{}, id).Error
}

func (s *AttachmentService) NormalizeURL(attachment *model.Attachment) {
	attachment.URL = "/attachments/" + strings.ReplaceAll(filepath.Base(attachment.Path), "\\", "/")
}
