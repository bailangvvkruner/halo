package extension

import (
	"fmt"
	"reflect"
	"sync"
)

var defaultScheme = NewScheme()

func DefaultScheme() *Scheme {
	return defaultScheme
}

type gvkKey struct {
	Group   string
	Version string
	Kind    string
}

func toGVKKey(gvk GVK) gvkKey {
	return gvkKey{Group: gvk.Group, Version: gvk.Version, Kind: gvk.Kind}
}

// Scheme 维护 GVK 到 Go 类型之间的映射注册表。
// 通过 Scheme，框架可以在运行时根据 GVK 信息实例化对应的 Go 结构体，
// 实现类似 Kubernetes Scheme 的类型发现和反序列化能力。
type Scheme struct {
	mu       sync.RWMutex
	gvkToType map[gvkKey]reflect.Type
	typeToGVK map[reflect.Type]GVK
}

// NewScheme 创建并返回一个新的 Scheme 实例。
func NewScheme() *Scheme {
	return &Scheme{
		gvkToType: make(map[gvkKey]reflect.Type),
		typeToGVK: make(map[reflect.Type]GVK),
	}
}

// Register 将一个 Go 类型注册到 Scheme 中，建立其与 GVK 的双向映射关系。
// obj 必须是一个指向实现了 Extension 接口的结构体指针（非 nil）。
// 如果同一 GVK 已被注册或类型已被映射到其他 GVK，将返回错误。
func (s *Scheme) Register(gvk GVK, obj Extension) error {
	if gvk.IsEmpty() {
		return fmt.Errorf("不能注册空的 GVK")
	}
	if obj == nil {
		return fmt.Errorf("注册对象不能为 nil")
	}

	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("注册对象必须是指针类型，实际为 %v", t.Kind())
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("注册对象必须指向结构体类型，实际底层为 %v", t.Kind())
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, exists := s.gvkToType[toGVKKey(gvk)]; exists {
		return fmt.Errorf("GVK %s 已被类型 %s 注册，不可重复注册", gvk, existing.Name())
	}
	if existingGVK, exists := s.typeToGVK[t]; exists {
		return fmt.Errorf("类型 %s 已注册到 GVK %s，不可重复注册", t.Name(), existingGVK)
	}

	s.gvkToType[toGVKKey(gvk)] = t
	s.typeToGVK[t] = gvk
	return nil
}

// KnownTypes 返回所有已注册的 GVK 列表。
func (s *Scheme) KnownTypes() []GVK {
	s.mu.RLock()
	defer s.mu.RUnlock()

	gvks := make([]GVK, 0, len(s.typeToGVK))
	for _, gvk := range s.typeToGVK {
		gvks = append(gvks, gvk)
	}
	return gvks
}

// New 根据 GVK 创建对应的扩展对象实例（零值）。
// 如果 GVK 未注册，返回错误。返回值是满足 Extension 接口的对象。
func (s *Scheme) New(gvk GVK) (Extension, error) {
	s.mu.RLock()
	t, ok := s.gvkToType[toGVKKey(gvk)]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("未找到 GVK %s 对应的已注册类型", gvk)
	}

	val := reflect.New(t).Interface()
	ext, ok := val.(Extension)
	if !ok {
		return nil, fmt.Errorf("GVK %s 对应的类型 %s 未实现 Extension 接口", gvk, t.Name())
	}
	return ext, nil
}

// MustNew 类似于 New，但在 GVK 未注册时触发 panic 而非返回错误。
// 适用于启动阶段确定所有类型已正确注册的场景。
func (s *Scheme) MustNew(gvk GVK) Extension {
	ext, err := s.New(gvk)
	if err != nil {
		panic(fmt.Sprintf("extension.Scheme.MustNew: %v", err))
	}
	return ext
}

// GVKFor 根据给定的 Extension 实例查询其对应的 GVK 信息。
// 如果该实例的类型未在 Scheme 中注册，返回 false。
func (s *Scheme) GVKFor(obj Extension) (GVK, bool) {
	if obj == nil {
		return GVK{}, false
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	s.mu.RLock()
	gvk, ok := s.typeToGVK[t]
	s.mu.RUnlock()
	return gvk, ok
}

// IsRegistered 判断指定的 GVK 是否已在 Scheme 中注册。
func (s *Scheme) IsRegistered(gvk GVK) bool {
	s.mu.RLock()
	_, ok := s.gvkToType[toGVKKey(gvk)]
	s.mu.RUnlock()
	return ok
}
