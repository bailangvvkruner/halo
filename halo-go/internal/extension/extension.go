package extension

// Extension 是所有扩展资源必须实现的接口。
// 通过此接口，框架可以统一处理不同类型的扩展对象，
// 获取其元数据和 GVK 标识信息，实现多态操作。
type Extension interface {
	// GetMetadata 返回扩展对象的元数据信息。
	GetMetadata() *Metadata

	// GetGVK 返回扩展对象的 GroupVersionKind 标识。
	GetGVK() GVK
}

// AbstractExtension 是扩展资源的抽象基础结构体，嵌入了 Metadata。
// 所有具体的扩展类型应匿名嵌入此结构体以自动满足 Extension 接口的部分要求。
// 具体类型仍需自行实现 GetGVK 方法以返回正确的 GVK 信息。
type AbstractExtension struct {
	// Metadata 包含该扩展对象的元数据信息（名称、标签、注解等）。
	Metadata Metadata `json:"metadata"`
}

// GetMetadata 返回嵌入的 Metadata 指针，满足 Extension 接口要求。
func (e *AbstractExtension) GetMetadata() *Metadata {
	return &e.Metadata
}
