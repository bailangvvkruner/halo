package extension

import "time"

// Metadata 包含所有扩展对象通用的元数据信息。
// 类似 Kubernetes ObjectMeta，作为扩展资源的主键和属性集合。
// 所有扩展对象必须嵌入此结构体以获得标准的元数据能力。
type Metadata struct {
	// Name 是扩展对象的唯一标识符（主键）。
	// 在同一 GVK 下全局唯一。
	Name string `json:"name"`

	// Generation 是一个单调递增的序列号，
	// 每次对象规格（Spec）变更时递增，用于观察状态变更。
	Generation int64 `json:"generation,omitempty"`

	// Labels 是键值对形式的标签集合，用于组织和筛选扩展对象。
	// 标签常用于查询选择器和分类管理。
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations 是非标识性的键值对元数据，用于存储工具或外部系统的附加信息。
	// 与 Labels 不同，Annotations 不用于对象选择。
	Annotations map[string]string `json:"annotations,omitempty"`

	// CreationTimestamp 表示对象的创建时间戳。
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty"`

	// DeletionTimestamp 表示对象的删除时间戳。
	// 当该字段不为 nil 时，表示对象正在被删除（优雅删除过程中）。
	DeletionTimestamp *time.Time `json:"deletionTimestamp,omitempty"`

	// Finalizers 是必须在对象删除前执行的清理逻辑列表。
	// 当 Finalizers 不为空时，即使收到删除请求也不会真正删除对象，
	// 必须等待所有 Finalizer 被移除后才会执行物理删除。
	Finalizers []string `json:"finalizers,omitempty"`

	// Version 用于乐观并发控制，每次更新时由服务端自动递增。
	// 客户端在更新时需携带当前 Version 以防止并发冲突。
	Version int64 `json:"version,omitempty"`
}

// NewMetadata 创建并返回一个新的 Metadata 实例，使用给定的名称初始化。
func NewMetadata(name string) *Metadata {
	return &Metadata{
		Name:              name,
		CreationTimestamp: time.Now(),
		Labels:            make(map[string]string),
		Annotations:       make(map[string]string),
	}
}

// SetLabel 设置指定键的标签值。
func (m *Metadata) SetLabel(key, value string) {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}
	m.Labels[key] = value
}

// GetLabel 获取指定键的标签值。如果不存在则返回空字符串。
func (m *Metadata) GetLabel(key string) string {
	if m.Labels == nil {
		return ""
	}
	return m.Labels[key]
}

// RemoveLabel 移除指定键的标签。
func (m *Metadata) RemoveLabel(key string) {
	delete(m.Labels, key)
}

// SetAnnotation 设置指定键的注解值。
func (m *Metadata) SetAnnotation(key, value string) {
	if m.Annotations == nil {
		m.Annotations = make(map[string]string)
	}
	m.Annotations[key] = value
}

// GetAnnotation 获取指定键的注解值。如果不存在则返回空字符串。
func (m *Metadata) GetAnnotation(key string) string {
	if m.Annotations == nil {
		return ""
	}
	return m.Annotations[key]
}

// HasFinalizer 判断指定的 Finalizer 是否存在于列表中。
func (m *Metadata) HasFinalizer(finalizer string) bool {
	for _, f := range m.Finalizers {
		if f == finalizer {
			return true
		}
	}
	return false
}

// AddFinalizer 向 Finalizers 列表中添加一个 Finalizer（去重）。
func (m *Metadata) AddFinalizer(finalizer string) {
	if m.HasFinalizer(finalizer) {
		return
	}
	m.Finalizers = append(m.Finalizers, finalizer)
}

// RemoveFinalizer 从 Finalizers 列表中移除指定的 Finalizer。
func (m *Metadata) RemoveFinalizer(finalizer string) {
	result := make([]string, 0, len(m.Finalizers))
	for _, f := range m.Finalizers {
		if f != finalizer {
			result = append(result, f)
		}
	}
	m.Finalizers = result
}

// IsDeleting 判断对象是否处于正在被删除的状态（即 DeletionTimestamp 不为 nil）。
func (m *Metadata) IsDeleting() bool {
	return m.DeletionTimestamp != nil
}
