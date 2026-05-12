package model

import "github.com/halo-dev/halo-go/internal/extension"

type TagSpec struct {
	DisplayName string `json:"displayName"`
	Slug        string `json:"slug"`
	Cover       string `json:"cover,omitempty"`
}

type TagStatus struct {
	Permalink string `json:"permalink"`
}

type Tag struct {
	extension.AbstractExtension
	Spec   TagSpec   `json:"spec"`
	Status TagStatus `json:"status,omitempty"`
}

func (t *Tag) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "content.halo.run",
		Version:  "v1alpha1",
		Kind:     "Tag",
		Plural:   "tags",
		Singular: "tag",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Tag{}).GetGVK(), &Tag{})
}
