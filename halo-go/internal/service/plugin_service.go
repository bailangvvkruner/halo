package service

import (
	"context"
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type PluginService interface {
	Create(ctx context.Context, plugin *model.Plugin) (*model.Plugin, error)
	Get(ctx context.Context, name string) (*model.Plugin, error)
	Update(ctx context.Context, plugin *model.Plugin) (*model.Plugin, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
	Enable(ctx context.Context, name string) error
	Disable(ctx context.Context, name string) error
}

type pluginService struct {
	client extension.TypedClient
}

func NewPluginService(client extension.TypedClient) PluginService {
	return &pluginService{client: client}
}

func (s *pluginService) Create(ctx context.Context, plugin *model.Plugin) (*model.Plugin, error) {
	if plugin.Metadata.Name == "" {
		plugin.Metadata.Name = generateName()
	}
	now := time.Now()
	plugin.Spec.SetupTime = &now
	if plugin.Spec.Enabled {
		plugin.Spec.StartTime = &now
	}
	ext, err := s.client.Create(ctx, plugin)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Plugin), nil
}

func (s *pluginService) Get(ctx context.Context, name string) (*model.Plugin, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Plugin), nil
}

func (s *pluginService) Update(ctx context.Context, plugin *model.Plugin) (*model.Plugin, error) {
	ext, err := s.client.Update(ctx, plugin)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Plugin), nil
}

func (s *pluginService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *pluginService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}

func (s *pluginService) Enable(ctx context.Context, name string) error {
	p, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	p.Spec.Enabled = true
	now := time.Now()
	p.Spec.StartTime = &now
	_, err = s.client.Update(ctx, p)
	return err
}

func (s *pluginService) Disable(ctx context.Context, name string) error {
	p, err := s.Get(ctx, name)
	if err != nil {
		return err
	}
	p.Spec.Enabled = false
	p.Spec.StartTime = nil
	_, err = s.client.Update(ctx, p)
	return err
}
