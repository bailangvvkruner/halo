package service

import (
	"halo/internal/model"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) List() ([]model.Post, error) {
	var posts []model.Post
	return posts, s.db.Order("id desc").Find(&posts).Error
}

func (s *PostService) FindBySlug(slug string) (*model.Post, error) {
	var post model.Post
	if err := s.db.Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) Create(post *model.Post) error {
	return s.db.Create(post).Error
}

func (s *PostService) Get(id uint) (*model.Post, error) {
	var post model.Post
	if err := s.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) Update(id uint, post *model.Post) error {
	return s.db.Model(&model.Post{}).Where("id = ?", id).Updates(post).Error
}

func (s *PostService) Delete(id uint) error {
	return s.db.Delete(&model.Post{}, id).Error
}

func (s *PostService) Search(keyword string, result *[]model.Post) error {
	return s.db.Where("title LIKE ? OR content LIKE ? OR excerpt LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Find(result).Error
}
