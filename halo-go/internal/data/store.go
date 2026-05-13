package data

import (
	"context"
	"encoding/json"
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
	"gorm.io/gorm"
)

type extensionRecord struct {
	ID        uint      `gorm:"primaryKey"`
	GroupName string    `gorm:"column:group_name;size:255;not null;index"`
	Version   string    `gorm:"column:version;size:64;not null;index"`
	Kind      string    `gorm:"column:kind;size:128;not null;index"`
	Name      string    `gorm:"column:name;size:255;not null;uniqueIndex"`
	Data      string    `gorm:"column:data;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (extensionRecord) TableName() string {
	return "extensions"
}

type ExtensionStore struct {
	db *gorm.DB
}

func NewExtensionStore(db *gorm.DB) *ExtensionStore {
	return &ExtensionStore{db: db}
}

func (s *ExtensionStore) Create(ctx context.Context, obj extension.Extension) (extension.Extension, error) {
	gvk := obj.GetGVK()
	meta := obj.GetMetadata()
	if meta.Name == "" {
		return nil, extension.ErrInvalidGVK
	}

	meta.CreationTimestamp = time.Now()
	meta.Version = 1
	meta.Generation = 1

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	record := extensionRecord{
		GroupName: gvk.Group,
		Version:   gvk.Version,
		Kind:      gvk.Kind,
		Name:      meta.Name,
		Data:      string(data),
		CreatedAt: meta.CreationTimestamp,
		UpdatedAt: meta.CreationTimestamp,
	}

	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		return nil, extension.ErrExtensionAlreadyExists
	}

	return obj, nil
}

func (s *ExtensionStore) Get(ctx context.Context, name string) (extension.Extension, error) {
	var record extensionRecord
	if err := s.db.WithContext(ctx).Where("name = ?", name).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, extension.ErrExtensionNotFound
		}
		return nil, err
	}

	gvk := extension.GVK{
		Group:   record.GroupName,
		Version: record.Version,
		Kind:    record.Kind,
	}

	ext, err := extension.DefaultScheme().New(gvk)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(record.Data), ext); err != nil {
		return nil, err
	}

	return ext, nil
}

func (s *ExtensionStore) Update(ctx context.Context, obj extension.Extension) (extension.Extension, error) {
	meta := obj.GetMetadata()
	var record extensionRecord
	if err := s.db.WithContext(ctx).Where("name = ?", meta.Name).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, extension.ErrExtensionNotFound
		}
		return nil, err
	}

	meta.Generation++
	meta.Version++
	now := time.Now()

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	record.Data = string(data)
	record.UpdatedAt = now

	if err := s.db.WithContext(ctx).Save(&record).Error; err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *ExtensionStore) Delete(ctx context.Context, name string) error {
	result := s.db.WithContext(ctx).Where("name = ?", name).Delete(&extensionRecord{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return extension.ErrExtensionNotFound
	}
	return nil
}

func (s *ExtensionStore) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	if opts == nil {
		opts = extension.DefaultListOptions()
	}
	opts.Validate()

	var total int64
	query := s.db.WithContext(ctx).Model(&extensionRecord{})
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var records []extensionRecord
	q := s.db.WithContext(ctx).Model(&extensionRecord{}).Order("created_at DESC")

	if opts.Size > 0 {
		q = q.Offset(opts.Offset()).Limit(opts.Size)
	}

	if err := q.Find(&records).Error; err != nil {
		return nil, err
	}

	items := make([]extension.Extension, 0, len(records))
	for _, record := range records {
		gvk := extension.GVK{
			Group:   record.GroupName,
			Version: record.Version,
			Kind:    record.Kind,
		}
		ext, err := extension.DefaultScheme().New(gvk)
		if err != nil {
			continue
		}
		if err := json.Unmarshal([]byte(record.Data), ext); err != nil {
			continue
		}
		items = append(items, ext)
	}

	return &extension.ListResult[extension.Extension]{
		Total: int(total),
		Page:  opts.Page,
		Size:  opts.Size,
		Items: items,
	}, nil
}

func (s *ExtensionStore) Watch(ctx context.Context, opts *extension.ListOptions) (<-chan extension.WatchEvent, error) {
	ch := make(chan extension.WatchEvent)
	go func() {
		<-ctx.Done()
		close(ch)
	}()
	return ch, nil
}

type typedClient struct {
	store *ExtensionStore
}

func NewClient(store *ExtensionStore) extension.TypedClient {
	return &typedClient{store: store}
}

func (c *typedClient) Create(ctx context.Context, obj extension.Extension) (extension.Extension, error) {
	return c.store.Create(ctx, obj)
}

func (c *typedClient) Get(ctx context.Context, name string) (extension.Extension, error) {
	return c.store.Get(ctx, name)
}

func (c *typedClient) Update(ctx context.Context, obj extension.Extension) (extension.Extension, error) {
	return c.store.Update(ctx, obj)
}

func (c *typedClient) Delete(ctx context.Context, name string) error {
	return c.store.Delete(ctx, name)
}

func (c *typedClient) List(ctx context.Context, opts *extension.ListOptions) (*extension.ListResult[extension.Extension], error) {
	return c.store.List(ctx, opts)
}

func (c *typedClient) Watch(ctx context.Context, opts *extension.ListOptions) (<-chan extension.WatchEvent, error) {
	return c.store.Watch(ctx, opts)
}