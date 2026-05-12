package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type MenuService interface {
	Create(ctx context.Context, menu *model.Menu) (*model.Menu, error)
	Get(ctx context.Context, name string) (*model.Menu, error)
	Update(ctx context.Context, menu *model.Menu) (*model.Menu, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
}

type menuService struct {
	client extension.TypedClient
}

func NewMenuService(client extension.TypedClient) MenuService {
	return &menuService{client: client}
}

func (s *menuService) Create(ctx context.Context, menu *model.Menu) (*model.Menu, error) {
	if menu.Metadata.Name == "" {
		menu.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, menu)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Menu), nil
}

func (s *menuService) Get(ctx context.Context, name string) (*model.Menu, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Menu), nil
}

func (s *menuService) Update(ctx context.Context, menu *model.Menu) (*model.Menu, error) {
	ext, err := s.client.Update(ctx, menu)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Menu), nil
}

func (s *menuService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *menuService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}
