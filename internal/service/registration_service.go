package service

import (
	"errors"
	"time"

	"halo/internal/model"

	"gorm.io/gorm"
)

type RegistrationService struct {
	db *gorm.DB
}

func NewRegistrationService(db *gorm.DB) *RegistrationService {
	return &RegistrationService{db: db}
}

func (s *RegistrationService) Register(reg *model.UserRegistration) error {
	var count int64
	if err := s.db.Model(&model.User{}).Where("username = ?", reg.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	reg.Token = GenerateToken(32)
	reg.ExpiresAt = time.Now().Add(24 * time.Hour)

	return s.db.Create(reg).Error
}

func (s *RegistrationService) Verify(token string) error {
	var reg model.UserRegistration
	if err := s.db.Where("token = ?", token).First(&reg).Error; err != nil {
		return err
	}

	if time.Now().After(reg.ExpiresAt) {
		return errors.New("token expired")
	}

	user := &model.User{
		Username: reg.Username,
		Password: reg.Password,
		Role:     "guest",
	}
	if err := s.db.Create(user).Error; err != nil {
		return err
	}

	reg.Verified = true
	s.db.Model(&reg).Update("verified", true)
	return nil
}

func (s *RegistrationService) GetPending(token string) (*model.UserRegistration, error) {
	var reg model.UserRegistration
	if err := s.db.Where("token = ?", token).First(&reg).Error; err != nil {
		return nil, err
	}
	return &reg, nil
}
