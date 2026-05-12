package service

import (
	"context"
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type PostService interface {
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	Get(ctx context.Context, name string) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) (*model.Post, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
	Publish(ctx context.Context, name string) error
	Unpublish(ctx context.Context, name string) error
	Trash(ctx context.Context, name string) error
	Restore(ctx context.Context, name string) error
}

type postService struct {
	client extension.TypedClient
	gvk    extension.GVK
}

func NewPostService(client extension.TypedClient) PostService {
	return &postService{
		client: client,
		gvk:    (&model.Post{}).GetGVK(),
	}
}

func (s *postService) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	if post.Metadata.Name == "" {
		post.Metadata.Name = generateName()
	}
	now := time.Now()
	if post.Spec.Publish && post.Spec.PublishTime == nil {
		post.Spec.PublishTime = &now
	}
	ext, err := s.client.Create(ctx, post)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Post), nil
}

func (s *postService) Get(ctx context.Context, name string) (*model.Post, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Post), nil
}

func (s *postService) Update(ctx context.Context, post *model.Post) (*model.Post, error) {
	ext, err := s.client.Update(ctx, post)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Post), nil
}

func (s *postService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *postService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}

func (s *postService) Publish(ctx context.Context, name string) error {
	post, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	now := time.Now()
	post.Spec.Publish = true
	post.Spec.PublishTime = &now
	_, err = s.client.Update(ctx, post)
	return err
}

func (s *postService) Unpublish(ctx context.Context, name string) error {
	post, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	post.Spec.Publish = false
	_, err = s.client.Update(ctx, post)
	return err
}

func (s *postService) Trash(ctx context.Context, name string) error {
	post, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	post.Spec.Deleted = true
	_, err = s.client.Update(ctx, post)
	return err
}

func (s *postService) Restore(ctx context.Context, name string) error {
	post, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	post.Spec.Deleted = false
	_, err = s.client.Update(ctx, post)
	return err
}
