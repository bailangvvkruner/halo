package extension

import (
	"fmt"
)

// Ref 表示对另一个扩展对象的引用，通过 GVK + Name 精确定位目标对象。
// 常用于表达扩展对象之间的关联关系（如插件依赖、父子关系等）。
type Ref struct {
	// Group 是被引用对象所属的 API 组名。
	Group string `json:"group"`

	// Version 是被引用对象的 API 版本号。
	Version string `json:"version"`

	// Kind 是被引用对象的资源类型名称。
	Kind string `json:"kind"`

	// Name 是被引用对象的唯一标识符（主键）。
	Name string `json:"name"`
}

// String 返回 Ref 的字符串表示形式：Group/Version/Kind/Name
func (r Ref) String() string {
	return fmt.Sprintf("%s/%s/%s/%s", r.Group, r.Version, r.Kind, r.Name)
}

// GVK 返回此引用对应的 GroupVersionKind 信息。
func (r Ref) GVK() GVK {
	return GVK{
		Group:   r.Group,
		Version: r.Version,
		Kind:    r.Kind,
	}
}

// Equal 判断两个 Ref 是否指向同一个目标对象。
func (r Ref) Equal(other Ref) bool {
	return r.Group == other.Group && r.Version == other.Version &&
		r.Kind == other.Kind && r.Name == other.Name
}

// IsEmpty 判断 Ref 是否为空值（所有字段均为空）。
func (r Ref) IsEmpty() bool {
	return r.Group == "" && r.Version == "" && r.Kind == "" && r.Name == ""
}

// NewRef 创建并返回一个新的 Ref 实例。
func NewRef(group, version, kind, name string) Ref {
	return Ref{
		Group:   group,
		Version: version,
		Kind:    kind,
		Name:    name,
	}
}

// NewRefFromExtension 从给定的 Extension 对象创建对应的 Ref 引用。
func NewRefFromExtension(ext Extension) Ref {
	gvk := ext.GetGVK()
	meta := ext.GetMetadata()
	return Ref{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
		Name:    meta.Name,
	}
}
