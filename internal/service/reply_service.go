package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type ReplyService struct {
	db *gorm.DB
}

func NewReplyService(db *gorm.DB) *ReplyService {
	return &ReplyService{db: db}
}

func (s *ReplyService) List() ([]model.Reply, error) {
	var replies []model.Reply
	return replies, s.db.Order("id desc").Find(&replies).Error
}

func (s *ReplyService) ListByComment(commentID uint) ([]model.Reply, error) {
	var replies []model.Reply
	return replies, s.db.Where("comment_id = ?", commentID).Order("id desc").Find(&replies).Error
}

func (s *ReplyService) Create(reply *model.Reply) error {
	reply.Status = "pending"
	return s.db.Create(reply).Error
}

func (s *ReplyService) Approve(id uint) error {
	return s.db.Model(&model.Reply{}).Where("id = ?", id).Update("status", "approved").Error
}

func (s *ReplyService) Reject(id uint) error {
	return s.db.Model(&model.Reply{}).Where("id = ?", id).Update("status", "rejected").Error
}

func (s *ReplyService) Delete(id uint) error {
	return s.db.Delete(&model.Reply{}, id).Error
}
