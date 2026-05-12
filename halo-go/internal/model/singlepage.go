package model

import "github.com/halo-dev/halo-go/internal/extension"

type SinglePageSpec struct {
	Title           string      `json:"title"`
	Slug            string      `json:"slug"`
	ReleaseSnapshot string      `json:"releaseSnapshot,omitempty"`
	HeadSnapshot    string      `json:"headSnapshot"`
	BaseSnapshot    string      `json:"baseSnapshot"`
	Owner           string      `json:"owner"`
	Template        string      `json:"template,omitempty"`
	Cover           string      `json:"cover,omitempty"`
	Deleted         bool        `json:"deleted"`
	Publish         bool        `json:"publish"`
	Pinned          bool        `json:"pinned"`
	AllowComment    bool        `json:"allowComment"`
	Visible         string      `json:"visible"`
	Version         int         `json:"version"`
	Priority        int         `json:"priority"`
	Excerpt         Excerpt     `json:"excerpt"`
	HtmlMetas       []HtmlMeta  `json:"htmlMetas,omitempty"`
}

type SinglePageStatus struct {
	Permalink string `json:"permalink"`
}

type SinglePage struct {
	extension.AbstractExtension
	Spec   SinglePageSpec   `json:"spec"`
	Status SinglePageStatus `json:"status,omitempty"`
}

func (sp *SinglePage) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "content.halo.run",
		Version:  "v1alpha1",
		Kind:     "SinglePage",
		Plural:   "singlepages",
		Singular: "singlepage",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&SinglePage{}).GetGVK(), &SinglePage{})
}
