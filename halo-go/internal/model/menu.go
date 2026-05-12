package model

import "github.com/halo-dev/halo-go/internal/extension"

type MenuSpec struct {
	DisplayName string   `json:"displayName"`
	MenuItems   []string `json:"menuItems,omitempty"`
}

type MenuStatus struct{}

type Menu struct {
	extension.AbstractExtension
	Spec   MenuSpec   `json:"spec"`
	Status MenuStatus `json:"status,omitempty"`
}

func (m *Menu) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "",
		Version:  "v1alpha1",
		Kind:     "Menu",
		Plural:   "menus",
		Singular: "menu",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Menu{}).GetGVK(), &Menu{})
}
