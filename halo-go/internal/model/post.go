package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type PostSpec struct {
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
	PublishTime     *time.Time  `json:"publishTime,omitempty"`
	Pinned          bool        `json:"pinned"`
	AllowComment    bool        `json:"allowComment"`
	Visible         string      `json:"visible"`
	Priority        int         `json:"priority"`
	Excerpt         Excerpt     `json:"excerpt"`
	Categories      []string    `json:"categories,omitempty"`
	Tags            []string    `json:"tags,omitempty"`
	HtmlMetas       []HtmlMeta  `json:"htmlMetas,omitempty"`
}

type Excerpt struct {
	AutoGenerate bool   `json:"autoGenerate"`
	Raw          string `json:"raw,omitempty"`
}

type HtmlMeta struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

type PostStatus struct {
	Permalink string `json:"permalink"`
}

type Post struct {
	extension.AbstractExtension
	Spec   PostSpec   `json:"spec"`
	Status PostStatus `json:"status,omitempty"`
}

func (p *Post) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "content.halo.run",
		Version:  "v1alpha1",
		Kind:     "Post",
		Plural:   "posts",
		Singular: "post",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Post{}).GetGVK(), &Post{})
}
