package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type NotificationReasonType string

const (
	NotificationReasonCommenterReply NotificationReasonType = "CommenterReply"
	NotificationReasonNewComment     NotificationReasonType = "NewComment"
)

type NotificationSpec struct {
	Title           string                  `json:"title"`
	Raw             string                  `json:"raw"`
	Content         string                  `json:"content"`
	Author          string                  `json:"author"`
	Receivers       []string                `json:"receivers"`
	UnreadReceivers []string                `json:"unreadReceivers"`
	Reason          NotificationReasonType  `json:"reason"`
	Attributes      map[string]string       `json:"attributes,omitempty"`
	PublishTime     *time.Time              `json:"publishTime,omitempty"`
}

type NotificationStatus struct{}

type Notification struct {
	extension.AbstractExtension
	Spec   NotificationSpec   `json:"spec"`
	Status NotificationStatus `json:"status,omitempty"`
}

func (n *Notification) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "notification.halo.run",
		Version:  "v1alpha1",
		Kind:     "Notification",
		Plural:   "notifications",
		Singular: "notification",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Notification{}).GetGVK(), &Notification{})
}
