package service

import (
	"errors"

	"halo/internal/model"

	"gorm.io/gorm"
)

type SetupService struct {
	db       *gorm.DB
	users    *UserService
	settings *SettingService
}

type SetupPayload struct {
	SiteTitle   string `json:"siteTitle"`
	BaseURL     string `json:"baseURL"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Language    string `json:"language"`
}

func NewSetupService(db *gorm.DB, users *UserService, settings *SettingService) *SetupService {
	return &SetupService{db: db, users: users, settings: settings}
}

func (s *SetupService) IsInitialized() (bool, error) {
	var count int64
	if err := s.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SetupService) Setup(payload SetupPayload) error {
	initialized, err := s.IsInitialized()
	if err != nil {
		return err
	}
	if initialized {
		return errors.New("system already initialized")
	}

	user := &model.User{
		Username: payload.Username,
		Password: payload.Password,
		Role:     "super-role",
	}
	if err := s.users.Create(user); err != nil {
		return err
	}

	settings := []model.Setting{
		{Key: "site.title", Value: payload.SiteTitle},
		{Key: "site.base_url", Value: payload.BaseURL},
		{Key: "site.language", Value: payload.Language},
		{Key: "site.admin_email", Value: payload.Email},
	}

	for _, item := range settings {
		copy := item
		if err := s.db.Create(&copy).Error; err != nil {
			return err
		}
	}

	return nil
}
