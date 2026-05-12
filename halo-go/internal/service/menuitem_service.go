package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type MenuItemService interface {
	Create(ctx context.Context, item *model.MenuItem) (*model.MenuItem, error)
	Get(ctx context.Context, name string) (*model.MenuItem, error)
	Update(ctx context.Context, item *model.MenuItem) (*model.MenuItem, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
}

type menuItemService struct {
	client extension.TypedClient
}

func NewMenuItemService(client extension.TypedClient) MenuItemService {
	return &menuItemService{client: client}
}

func (s *menuItemService) Create(ctx context.Context, item *model.MenuItem) (*model.MenuItem, error) {
	if item.Metadata.Name == "" {
		item.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return ext.(*model.MenuItem), nil
}

func (s *menuItemService) Get(ctx context.Context, name string) (*model.MenuItem, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.MenuItem), nil
}

func (s *menuItemService) Update(ctx context.Context, item *model.MenuItem) (*model.MenuItem, error) {
	ext, err := s.client.Update(ctx, item)
	if err != nil {
		return nil, err
	}
	return ext.(*model.MenuItem), nil
}

func (s *menuItemService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *menuItemService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}
