package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) List() ([]model.Comment, error) {
	var comments []model.Comment
	return comments, s.db.Order("id desc").Find(&comments).Error
}

func (s *CommentService) Create(comment *model.Comment) error {
	return s.db.Create(comment).Error
}

func (s *CommentService) Get(id uint) (*model.Comment, error) {
	var item model.Comment
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *CommentService) Delete(id uint) error {
	return s.db.Delete(&model.Comment{}, id).Error
}
