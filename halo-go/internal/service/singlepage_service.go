package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type SinglePageService interface {
	Create(ctx context.Context, sp *model.SinglePage) (*model.SinglePage, error)
	Get(ctx context.Context, name string) (*model.SinglePage, error)
	Update(ctx context.Context, sp *model.SinglePage) (*model.SinglePage, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
	Publish(ctx context.Context, name string) error
	Unpublish(ctx context.Context, name string) error
	Trash(ctx context.Context, name string) error
	Restore(ctx context.Context, name string) error
}

type singlePageService struct {
	client extension.TypedClient
	gvk    extension.GVK
}

func NewSinglePageService(client extension.TypedClient) SinglePageService {
	return &singlePageService{
		client: client,
		gvk:    (&model.SinglePage{}).GetGVK(),
	}
}

func (s *singlePageService) Create(ctx context.Context, sp *model.SinglePage) (*model.SinglePage, error) {
	if sp.Metadata.Name == "" {
		sp.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, sp)
	if err != nil {
		return nil, err
	}
	return ext.(*model.SinglePage), nil
}

func (s *singlePageService) Get(ctx context.Context, name string) (*model.SinglePage, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.SinglePage), nil
}

func (s *singlePageService) Update(ctx context.Context, sp *model.SinglePage) (*model.SinglePage, error) {
	ext, err := s.client.Update(ctx, sp)
	if err != nil {
		return nil, err
	}
	return ext.(*model.SinglePage), nil
}

func (s *singlePageService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *singlePageService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}

func (s *singlePageService) Publish(ctx context.Context, name string) error {
	sp, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	sp.Spec.Publish = true
	_, err = s.client.Update(ctx, sp)
	return err
}

func (s *singlePageService) Unpublish(ctx context.Context, name string) error {
	sp, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	sp.Spec.Publish = false
	_, err = s.client.Update(ctx, sp)
	return err
}

func (s *singlePageService) Trash(ctx context.Context, name string) error {
	sp, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	sp.Spec.Deleted = true
	_, err = s.client.Update(ctx, sp)
	return err
}

func (s *singlePageService) Restore(ctx context.Context, name string) error {
	sp, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	sp.Spec.Deleted = false
	_, err = s.client.Update(ctx, sp)
	return err
}
