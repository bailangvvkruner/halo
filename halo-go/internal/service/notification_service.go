package service

import (
	"context"
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type NotificationService interface {
	Create(ctx context.Context, notification *model.Notification) (*model.Notification, error)
	Get(ctx context.Context, name string) (*model.Notification, error)
	Update(ctx context.Context, notification *model.Notification) (*model.Notification, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
	MarkRead(ctx context.Context, name string, receiver string) error
}

type notificationService struct {
	client extension.TypedClient
}

func NewNotificationService(client extension.TypedClient) NotificationService {
	return &notificationService{client: client}
}

func (s *notificationService) Create(ctx context.Context, notification *model.Notification) (*model.Notification, error) {
	if notification.Metadata.Name == "" {
		notification.Metadata.Name = generateName()
	}
	now := time.Now()
	notification.Spec.PublishTime = &now
	ext, err := s.client.Create(ctx, notification)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Notification), nil
}

func (s *notificationService) Get(ctx context.Context, name string) (*model.Notification, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Notification), nil
}

func (s *notificationService) Update(ctx context.Context, notification *model.Notification) (*model.Notification, error) {
	ext, err := s.client.Update(ctx, notification)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Notification), nil
}

func (s *notificationService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *notificationService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}

func (s *notificationService) MarkRead(ctx context.Context, name string, receiver string) error {
	n, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	var updated []string
	for _, r := range n.Spec.UnreadReceivers {
		if r != receiver {
			updated = append(updated, r)
		}
	}
	n.Spec.UnreadReceivers = updated
	_, err = s.client.Update(ctx, n)
	return err
}
