package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type ReplySpec struct {
	Raw            string     `json:"raw"`
	Content        string     `json:"content"`
	RawContent     string     `json:"rawContent"`
	Owner          string     `json:"owner"`
	CommentName    string     `json:"commentName"`
	QuoteReplyName string     `json:"quoteReplyName,omitempty"`
	AllowNotification bool   `json:"allowNotification"`
	Approved       bool       `json:"approved"`
	Reason         string     `json:"reason,omitempty"`
	ApproveTime    *time.Time `json:"approveTime,omitempty"`
	CreateTime     time.Time  `json:"createTime"`
	LastModifyTime time.Time  `json:"lastModifyTime"`
}

type ReplyStatus struct{}

type Reply struct {
	extension.AbstractExtension
	Spec   ReplySpec   `json:"spec"`
	Status ReplyStatus `json:"status,omitempty"`
}

func (r *Reply) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "content.halo.run",
		Version:  "v1alpha1",
		Kind:     "Reply",
		Plural:   "replies",
		Singular: "reply",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&Reply{}).GetGVK(), &Reply{})
}
