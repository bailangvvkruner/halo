package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type ThemeService interface {
	Create(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	Get(ctx context.Context, name string) (*model.Theme, error)
	Update(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
	Activate(ctx context.Context, name string) error
}

type themeService struct {
	client extension.TypedClient
}

func NewThemeService(client extension.TypedClient) ThemeService {
	return &themeService{client: client}
}

func (s *themeService) Create(ctx context.Context, theme *model.Theme) (*model.Theme, error) {
	if theme.Metadata.Name == "" {
		theme.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, theme)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Theme), nil
}

func (s *themeService) Get(ctx context.Context, name string) (*model.Theme, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Theme), nil
}

func (s *themeService) Update(ctx context.Context, theme *model.Theme) (*model.Theme, error) {
	ext, err := s.client.Update(ctx, theme)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Theme), nil
}

func (s *themeService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *themeService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}

func (s *themeService) Activate(ctx context.Context, name string) error {
	result, err := s.client.List(ctx, &extension.ListOptions{Size: 0})
	if err != nil {
		return err
	}
	for _, ext := range result.Items {
		t := ext.(*model.Theme)
		t.Spec.Active = t.Metadata.Name == name
		_, err := s.client.Update(ctx, t)
		if err != nil {
			return err
		}
	}
	return nil
}
