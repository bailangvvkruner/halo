package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type CategoryService interface {
	Create(ctx context.Context, category *model.Category) (*model.Category, error)
	Get(ctx context.Context, name string) (*model.Category, error)
	Update(ctx context.Context, category *model.Category) (*model.Category, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
}

type categoryService struct {
	client extension.TypedClient
}

func NewCategoryService(client extension.TypedClient) CategoryService {
	return &categoryService{client: client}
}

func (s *categoryService) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	if category.Metadata.Name == "" {
		category.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, category)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Category), nil
}

func (s *categoryService) Get(ctx context.Context, name string) (*model.Category, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Category), nil
}

func (s *categoryService) Update(ctx context.Context, category *model.Category) (*model.Category, error) {
	ext, err := s.client.Update(ctx, category)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Category), nil
}

func (s *categoryService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *categoryService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}
