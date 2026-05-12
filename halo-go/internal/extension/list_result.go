package extension

import "math"

// ListResult 是分页查询结果的泛型封装。
// 包含分页元数据和当前页的数据项列表，适用于所有扩展资源的列表查询返回值。
type ListResult[T any] struct {
	// Total 表示符合查询条件的总记录数（不受分页限制）。
	Total int `json:"total"`

	// Page 表示当前页码（从 1 开始）。
	Page int `json:"page"`

	// Size 表示每页的记录条数。
	Size int `json:"size"`

	// Items 是当前页的扩展对象列表。
	Items []T `json:"items"`
}

// First 返回第一页的页码（固定为 1）。
func (r *ListResult[T]) First() int {
	return 1
}

// Last 返回最后一页的页码。
// 当没有数据时返回 0，当 Size 为 0 时返回 1 以避免除零错误。
func (r *ListResult[T]) Last() int {
	if r.Total == 0 {
		return 0
	}
	if r.Size <= 0 {
		return 1
	}
	return int(math.Ceil(float64(r.Total) / float64(r.Size)))
}

// HasNext 判断是否存在下一页。
func (r *ListResult[T]) HasNext() bool {
	return r.Page < r.Last()
}

// HasPrev 判断是否存在上一页。
func (r *ListResult[T]) HasPrev() bool {
	return r.Page > 1
}

// TotalPages 返回总页数。
func (r *ListResult[T]) TotalPages() int {
	return r.Last()
}

// IsEmpty 判断当前页是否为空（无数据项）。
func (r *ListResult[T]) IsEmpty() bool {
	return len(r.Items) == 0
}

// Offset 返回当前页在全局结果集中的偏移量（从 0 开始）。
func (r *ListResult[T]) Offset() int {
	if r.Page <= 0 || r.Size <= 0 {
		return 0
	}
	return (r.Page - 1) * r.Size
}

// NewListResult 创建并返回一个新的 ListResult 实例。
func NewListResult[T any](total, page, size int, items []T) *ListResult[T] {
	return &ListResult[T]{
		Total: total,
		Page:  page,
		Size:  size,
		Items: items,
	}
}
