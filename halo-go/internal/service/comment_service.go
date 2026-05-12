package service

import (
	"context"
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type CommentService interface {
	Create(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	Get(ctx context.Context, name string) (*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
	Approve(ctx context.Context, name string) error
	ReplyTo(ctx context.Context, commentName string, reply *model.Reply) error
}

type commentService struct {
	client   extension.TypedClient
	replySvc ReplyService
}

func NewCommentService(client extension.TypedClient, replySvc ReplyService) CommentService {
	return &commentService{client: client, replySvc: replySvc}
}

func (s *commentService) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	if comment.Metadata.Name == "" {
		comment.Metadata.Name = generateName()
	}
	comment.Spec.CreateTime = time.Now()
	comment.Spec.LastModifyTime = time.Now()
	ext, err := s.client.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Comment), nil
}

func (s *commentService) Get(ctx context.Context, name string) (*model.Comment, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Comment), nil
}

func (s *commentService) Update(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	comment.Spec.LastModifyTime = time.Now()
	ext, err := s.client.Update(ctx, comment)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Comment), nil
}

func (s *commentService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *commentService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}

func (s *commentService) Approve(ctx context.Context, name string) error {
	comment, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	comment.Spec.Approved = true
	_, err = s.client.Update(ctx, comment)
	return err
}

func (s *commentService) ReplyTo(ctx context.Context, commentName string, reply *model.Reply) error {
	reply.Spec.CommentName = commentName
	_, err := s.replySvc.Create(ctx, reply)
	return err
}
