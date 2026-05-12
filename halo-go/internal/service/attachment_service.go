package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type AttachmentService interface {
	Create(ctx context.Context, attachment *model.Attachment) (*model.Attachment, error)
	Get(ctx context.Context, name string) (*model.Attachment, error)
	Update(ctx context.Context, attachment *model.Attachment) (*model.Attachment, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
}

type attachmentService struct {
	client extension.TypedClient
}

func NewAttachmentService(client extension.TypedClient) AttachmentService {
	return &attachmentService{client: client}
}

func (s *attachmentService) Create(ctx context.Context, attachment *model.Attachment) (*model.Attachment, error) {
	if attachment.Metadata.Name == "" {
		attachment.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, attachment)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Attachment), nil
}

func (s *attachmentService) Get(ctx context.Context, name string) (*model.Attachment, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Attachment), nil
}

func (s *attachmentService) Update(ctx context.Context, attachment *model.Attachment) (*model.Attachment, error) {
	ext, err := s.client.Update(ctx, attachment)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Attachment), nil
}

func (s *attachmentService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *attachmentService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}
