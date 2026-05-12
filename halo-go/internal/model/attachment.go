package model

import "github.com/halo-dev/halo-go/internal/extension"

type AttachmentSpec struct {
	DisplayName string `json:"displayName"`
	MediaType   string `json:"mediaType"`
	Size        int64  `json:"size"`
	Owner       string `json:"owner"`
	Path        string `json:"path"`
}

type AttachmentStatus struct {
	Permalink string `json:"permalink"`
}

type Attachment struct {
	extension.AbstractExtension
	Spec   AttachmentSpec   `json:"spec"`
	Status AttachmentStatus `json:"status,omitempty"`
}

func (a *Attachment) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "storage.halo.run",
		Version:  "v1alpha1",
		Kind:     "Attachment",
		Plural:   "attachments",
		Singular: "attachment",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Attachment{}).GetGVK(), &Attachment{})
}
