package model

import (
	"github.com/halo-dev/halo-go/internal/extension"
)

type CounterSpec struct {
	VisitTotal        int64 `json:"visitTotal"`
	UniqueVisitor     int64 `json:"uniqueVisitor"`
	PostVisitTotal    int64 `json:"postVisitTotal"`
	PostUniqueVisitor int64 `json:"postUniqueVisitor"`
}

type Counter struct {
	extension.AbstractExtension `json:",inline"`
	Spec                        CounterSpec `json:"spec"`
}

func newCounter() *Counter { return &Counter{} }

func (c *Counter) GetGVK() extension.GVK {
	return extension.NewGVK("", "v1alpha1", "Counter", "counters", "counter")
}
