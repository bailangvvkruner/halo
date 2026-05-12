package extension

import (
	"context"
	"errors"
)

// 定义扩展客户端操作中可能出现的通用错误。
var (
	// ErrExtensionNotFound 表示请求的扩展对象不存在。
	ErrExtensionNotFound = errors.New("扩展对象不存在")

	// ErrExtensionAlreadyExists 表示要创建的扩展对象已存在（名称冲突）。
	ErrExtensionAlreadyExists = errors.New("扩展对象已存在")

	// ErrVersionConflict 表示更新时版本号冲突（乐观并发控制失败）。
	ErrVersionConflict = errors.New("版本冲突，请刷新后重试")

	// ErrInvalidGVK 表示 GVK 无效或未注册。
	ErrInvalidGVK = errors.New("无效的 GroupVersionKind")
)

// WatchEventType 表示变更事件类型的枚举。
type WatchEventType string

const (
	// WatchEventAdded 表示新增事件。
	WatchEventAdded WatchEventType = "ADDED"
	// WatchEventModified 表示修改事件。
	WatchEventModified WatchEventType = "MODIFIED"
	// WatchEventDeleted 表示删除事件。
	WatchEventDeleted WatchEventType = "DELETED"
	// WatchEventBookmark 表示书签事件（用于断点续传）。
	WatchEventBookmark WatchEventType = "BOOKMARK"
	// WatchEventError 表示错误事件。
	WatchEventError WatchEventType = "ERROR"
)

// WatchEvent 表示一次资源变更事件，包含事件类型和对应的扩展对象。
type WatchEvent struct {
	// Type 是此次事件的类型（新增/修改/删除等）。
	Type WatchEventType `json:"type"`

	// Object 是事件关联的扩展对象。对于删除事件，该对象可能仅包含元数据。
	Object Extension `json:"object"`
}

// ExtensionClient 定义了扩展资源的 CRUD 和监听操作的统一接口。
// 所有存储后端（SQLite、内存、远程等）均需实现此接口，
// 以确保上层业务逻辑与具体存储实现解耦。
//
// 类型参数 T 约束了此客户端操作的扩展类型，必须满足 Extension 接口。
type ExtensionClient[T Extension] interface {
	// Create 创建一个新的扩展对象并持久化到存储中。
	// 如果同名对象已存在，应返回 ErrExtensionAlreadyExists 错误。
	Create(ctx context.Context, obj T) (T, error)

	// Get 根据名称获取单个扩展对象。
	// 如果对象不存在，应返回 ErrExtensionNotFound 错误。
	Get(ctx context.Context, name string) (T, error)

	// Update 更新一个已存在的扩展对象。
	// 实现方应根据 Metadata.Version 进行乐观并发控制：
	//   - 当传入对象的 Version 与存储中的当前值不匹配时，返回 ErrVersionConflict。
	//   - 更新成功后，应自动递增 Version 并更新 Generation。
	Update(ctx context.Context, obj T) (T, error)

	// Delete 按名称删除指定的扩展对象。
	// 如果对象存在 Finalizers，则不应立即物理删除，而是设置 DeletionTimestamp
	// 进入优雅删除状态；当所有 Finalizer 被移除后再执行物理删除。
	Delete(ctx context.Context, name string) error

	// List 根据查询选项分页列出符合条件的扩展对象。
	List(ctx context.Context, opts *ListOptions) (*ListResult[T], error)

	// Watch 监听指定 GVK 的所有资源变更事件。
	// 返回一个只读的事件通道，调用方通过消费通道来接收实时变更通知。
	// ctx 用于控制监听的生命周期，取消 context 将关闭事件通道。
	Watch(ctx context.Context, opts *ListOptions) (<-chan WatchEvent, error)
}

// TypedClient 是非泛型的扩展客户端接口变体。
// 适用于需要运行时动态确定类型或无法使用泛型的场景（如反射调用、路由分发）。
type TypedClient interface {
	// Create 创建扩展对象（使用 Extension 接口类型）。
	Create(ctx context.Context, obj Extension) (Extension, error)

	// Get 根据名称获取扩展对象。
	Get(ctx context.Context, name string) (Extension, error)

	// Update 更新扩展对象。
	Update(ctx context.Context, obj Extension) (Extension, error)

	// Delete 删除扩展对象。
	Delete(ctx context.Context, name string) error

	// List 列出扩展对象（返回 Extension 接口切片）。
	List(ctx context.Context, opts *ListOptions) (*ListResult[Extension], error)

	// Watch 监听扩展对象变更事件。
	Watch(ctx context.Context, opts *ListOptions) (<-chan WatchEvent, error)
}
