package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type TagService interface {
	Create(ctx context.Context, tag *model.Tag) (*model.Tag, error)
	Get(ctx context.Context, name string) (*model.Tag, error)
	Update(ctx context.Context, tag *model.Tag) (*model.Tag, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
}

type tagService struct {
	client extension.TypedClient
}

func NewTagService(client extension.TypedClient) TagService {
	return &tagService{client: client}
}

func (s *tagService) Create(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	if tag.Metadata.Name == "" {
		tag.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, tag)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Tag), nil
}

func (s *tagService) Get(ctx context.Context, name string) (*model.Tag, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Tag), nil
}

func (s *tagService) Update(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	ext, err := s.client.Update(ctx, tag)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Tag), nil
}

func (s *tagService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *tagService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}
