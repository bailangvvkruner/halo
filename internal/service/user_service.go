package service

import (
	"errors"

	"halo/internal/model"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) EnsureAdmin() error {
	var count int64
	if err := s.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return s.db.Create(&model.User{Username: "admin", Password: "admin123", Role: "super-role"}).Error
}

func (s *UserService) FindAll() ([]model.User, error) {
	var users []model.User
	return users, s.db.Order("id desc").Find(&users).Error
}

func (s *UserService) Authenticate(username, password string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}
