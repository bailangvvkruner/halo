package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type UserService interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Get(ctx context.Context, name string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
	GetByName(ctx context.Context, userName string) (*model.User, error)
	ChangePassword(ctx context.Context, name string, oldPwd, newPwd string) error
	UpdateProfile(ctx context.Context, user *model.User) error
}

type userService struct {
	client extension.TypedClient
}

func NewUserService(client extension.TypedClient) UserService {
	return &userService{client: client}
}

func (s *userService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if user.Metadata.Name == "" {
		user.Metadata.Name = generateName()
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Spec.RawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}
	user.Spec.Password = string(hashedPwd)
	user.Spec.RawPassword = ""
	ext, err := s.client.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	u := ext.(*model.User)
	u.Spec.Password = ""
	return u, nil
}

func (s *userService) Get(ctx context.Context, name string) (*model.User, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	u := ext.(*model.User)
	u.Spec.Password = ""
	return u, nil
}

func (s *userService) Update(ctx context.Context, user *model.User) (*model.User, error) {
	user.Spec.RawPassword = ""
	user.Spec.Password = ""
	ext, err := s.client.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	u := ext.(*model.User)
	u.Spec.Password = ""
	return u, nil
}

func (s *userService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *userService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}

func (s *userService) GetByName(ctx context.Context, userName string) (*model.User, error) {
	result, err := s.client.List(ctx, &extension.ListOptions{Size: 0})
	if err != nil {
		return nil, err
	}
	var items = result.Items
	for i := 0; i < len(items); i++ {
		ext := items[i]
		u, ok := ext.(*model.User)
		if !ok {
			continue
		}
		if u.Spec.UserName == userName {
			return u, nil
		}
	}
	return nil, fmt.Errorf("用户不存在: %s", userName)
}

func (s *userService) ChangePassword(ctx context.Context, name string, oldPwd, newPwd string) error {
	user, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Spec.Password), []byte(oldPwd)); err != nil {
		return fmt.Errorf("旧密码不正确")
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("新密码加密失败: %w", err)
	}
	user.Spec.Password = string(hashedPwd)
	_, err = s.client.Update(ctx, user)
	return err
}

func (s *userService) UpdateProfile(ctx context.Context, user *model.User) error {
	existing, err := s.Get(ctx, user.Metadata.Name)
	if err != nil {
		return err
	}
	existing.Spec.DisplayName = user.Spec.DisplayName
	existing.Spec.Avatar = user.Spec.Avatar
	existing.Spec.Bio = user.Spec.Bio
	existing.Spec.Email = user.Spec.Email
	now := time.Now()
	existing.Spec.LastLoginTime = &now
	_, err = s.client.Update(ctx, existing)
	return err
}
