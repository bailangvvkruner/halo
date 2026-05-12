package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type SnapshotSpec struct {
	SubjectRef     extension.Ref `json:"subjectRef"`
	RawType        string        `json:"rawType"`
	RawPatch       string        `json:"rawPatch"`
	ContentPatch   string        `json:"contentPatch"`
	LastModifyTime time.Time     `json:"lastModifyTime"`
	Owner          string        `json:"owner"`
	Contributors   []string      `json:"contributors"`
}

type SnapshotStatus struct{}

type Snapshot struct {
	extension.AbstractExtension
	Spec   SnapshotSpec   `json:"spec"`
	Status SnapshotStatus `json:"status,omitempty"`
}

func (s *Snapshot) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "content.halo.run",
		Version:  "v1alpha1",
		Kind:     "Snapshot",
		Plural:   "snapshots",
		Singular: "snapshot",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Snapshot{}).GetGVK(), &Snapshot{})
}
