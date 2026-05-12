package service

import (
	"context"

	"github.com/halo-dev/halo-go/internal/extension"
	"github.com/halo-dev/halo-go/internal/model"
)

type SnapshotService interface {
	Create(ctx context.Context, snapshot *model.Snapshot) (*model.Snapshot, error)
	Get(ctx context.Context, name string) (*model.Snapshot, error)
	Update(ctx context.Context, snapshot *model.Snapshot) (*model.Snapshot, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error)
	ListBySubject(ctx context.Context, subjectRef extension.Ref) ([]*model.Snapshot, error)
}

type snapshotService struct {
	client extension.TypedClient
	gvk    extension.GVK
}

func NewSnapshotService(client extension.TypedClient) SnapshotService {
	return &snapshotService{
		client: client,
		gvk:    (&model.Snapshot{}).GetGVK(),
	}
}

func (s *snapshotService) Create(ctx context.Context, snapshot *model.Snapshot) (*model.Snapshot, error) {
	if snapshot.Metadata.Name == "" {
		snapshot.Metadata.Name = generateName()
	}
	ext, err := s.client.Create(ctx, snapshot)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Snapshot), nil
}

func (s *snapshotService) Get(ctx context.Context, name string) (*model.Snapshot, error) {
	ext, err := s.client.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Snapshot), nil
}

func (s *snapshotService) Update(ctx context.Context, snapshot *model.Snapshot) (*model.Snapshot, error) {
	ext, err := s.client.Update(ctx, snapshot)
	if err != nil {
		return nil, err
	}
	return ext.(*model.Snapshot), nil
}

func (s *snapshotService) Delete(ctx context.Context, name string) error {
	return s.client.Delete(ctx, name)
}

func (s *snapshotService) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return s.client.List(ctx, opts)
}

func (s *snapshotService) ListBySubject(ctx context.Context, subjectRef extension.Ref) ([]*model.Snapshot, error) {
	result, err := s.client.List(ctx, &extension.ListOptions{Size: 0})
	if err != nil {
		return nil, err
	}
	var snapshots []*model.Snapshot
	for _, ext := range result.Items {
		sn := ext.(*model.Snapshot)
		if sn.Spec.SubjectRef.Group == subjectRef.Group &&
			sn.Spec.SubjectRef.Version == subjectRef.Version &&
			sn.Spec.SubjectRef.Kind == subjectRef.Kind &&
			sn.Spec.SubjectRef.Name == subjectRef.Name {
			snapshots = append(snapshots, sn)
		}
	}
	return snapshots, nil
}
