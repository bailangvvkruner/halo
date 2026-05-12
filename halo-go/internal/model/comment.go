package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type CommentSpec struct {
	Raw              string           `json:"raw"`
	Content          string           `json:"content"`
	Owner            string           `json:"owner"`
	SubjectRef       extension.Ref    `json:"subjectRef"`
	AllowNotification bool            `json:"allowNotification"`
	Approved         bool             `json:"approved"`
	TopPriority      bool             `json:"topPriority"`
	CreateTime       time.Time        `json:"createTime"`
	LastModifyTime   time.Time        `json:"lastModifyTime"`
}

type CommentStatus struct {
	ReplyCount int64 `json:"replyCount"`
}

type Comment struct {
	extension.AbstractExtension
	Spec   CommentSpec   `json:"spec"`
	Status CommentStatus `json:"status,omitempty"`
}

func (c *Comment) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "content.halo.run",
		Version:  "v1alpha1",
		Kind:     "Comment",
		Plural:   "comments",
		Singular: "comment",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&Comment{}).GetGVK(), &Comment{})
}
