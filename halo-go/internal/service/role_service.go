package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type RoleService interface {
	Create(ctx context.Context, role *model.Role) (*model.Role, error)
	Get(ctx context.Context, name string) (*model.Role, error)
	Update(ctx context.Context, role *model.Role) (*model.Role, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
}

type roleService struct {
	client extension.TypedClient
}

func NewRoleService(client extension.TypedClient) RoleService {
	return &roleService{client: client}
}

func (s *roleService) Create(ctx context.Context, role *model.Role) (*model.Role, error) {
	if role.Metadata.Name == "" {
		role.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, role)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Role), nil
}

func (s *roleService) Get(ctx context.Context, name string) (*model.Role, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Role), nil
}

func (s *roleService) Update(ctx context.Context, role *model.Role) (*model.Role, error) {
	ext, err := s.client.Update(ctx, role)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Role), nil
}

func (s *roleService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *roleService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}
