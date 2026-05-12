package model

import "github.com/halo-dev/halo-go/internal/extension"

type CategorySpec struct {
	DisplayName string   `json:"displayName"`
	Slug        string   `json:"slug"`
	Description string   `json:"description,omitempty"`
	Cover       string   `json:"cover,omitempty"`
	Template    string   `json:"template,omitempty"`
	Priority    int      `json:"priority"`
	Children    []string `json:"children,omitempty"`
}

type CategoryStatus struct {
	Permalink string `json:"permalink"`
}

type Category struct {
	extension.AbstractExtension
	Spec   CategorySpec   `json:"spec"`
	Status CategoryStatus `json:"status,omitempty"`
}

func (c *Category) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "content.halo.run",
		Version:  "v1alpha1",
		Kind:     "Category",
		Plural:   "categories",
		Singular: "category",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Category{}).GetGVK(), &Category{})
}
