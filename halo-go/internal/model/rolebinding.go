package model

import "github.com/halo-dev/halo-go/internal/extension"

type RoleBindingSpec struct {
	SubjectRef extension.Ref `json:"subjectRef"`
	RoleRef    extension.Ref `json:"roleRef"`
}

type RoleBindingStatus struct{}

type RoleBinding struct {
	extension.AbstractExtension
	Spec   RoleBindingSpec   `json:"spec"`
	Status RoleBindingStatus `json:"status,omitempty"`
}

func (rb *RoleBinding) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "",
		Version:  "v1alpha1",
		Kind:     "RoleBinding",
		Plural:   "rolebindings",
		Singular: "rolebinding",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&RoleBinding{}).GetGVK(), &RoleBinding{})
}
