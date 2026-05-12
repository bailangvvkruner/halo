package service

import (
	"context"
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type ReplyService interface {
	Create(ctx context.Context, reply *model.Reply) (*model.Reply, error)
	Get(ctx context.Context, name string) (*model.Reply, error)
	Update(ctx context.Context, reply *model.Reply) (*model.Reply, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
}

type replyService struct {
	client extension.TypedClient
}

func NewReplyService(client extension.TypedClient) ReplyService {
	return &replyService{client: client}
}

func (s *replyService) Create(ctx context.Context, reply *model.Reply) (*model.Reply, error) {
	if reply.Metadata.Name == "" {
		reply.Metadata.Name = generateName()
	}
	reply.Spec.CreateTime = time.Now()
	reply.Spec.LastModifyTime = time.Now()
	ext, err := s.client.Create(ctx, reply)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Reply), nil
}

func (s *replyService) Get(ctx context.Context, name string) (*model.Reply, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Reply), nil
}

func (s *replyService) Update(ctx context.Context, reply *model.Reply) (*model.Reply, error) {
	reply.Spec.LastModifyTime = time.Now()
	ext, err := s.client.Update(ctx, reply)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Reply), nil
}

func (s *replyService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *replyService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}
