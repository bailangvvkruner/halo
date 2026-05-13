package model

import "github.com/halo-dev/halo-go/internal/extension"

type Rule struct {
	APIGroups      []string `json:"apiGroups,omitempty"`
	Resources      []string `json:"resources,omitempty"`
	ResourceNames  []string `json:"resourceNames,omitempty"`
	Verbs           []string `json:"verbs"`
	NonResourceURLs []string `json:"nonResourceURLs,omitempty"`
}

type RoleSpec struct {
	DisplayName string            `json:"displayName"`
	Type        string            `json:"type"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Rules       []Rule             `json:"rules"`
}

type RoleStatus struct{}

type Role struct {
	extension.AbstractExtension
	Spec   RoleSpec   `json:"spec"`
	Status RoleStatus `json:"status,omitempty"`
}

func (r *Role) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "",
		Version:  "v1alpha1",
		Kind:     "Role",
		Plural:   "roles",
		Singular: "role",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&Role{}).GetGVK(), &Role{})
}
