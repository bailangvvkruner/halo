package model

import "github.com/halo-dev/halo-go/internal/extension"

type MenuItemSpec struct {
	DisplayName string          `json:"displayName"`
	Href        string          `json:"href"`
	TargetRef   *extension.Ref  `json:"targetRef,omitempty"`
	Children    []string        `json:"children,omitempty"`
	Priority    int             `json:"priority"`
}

type MenuItemStatus struct{}

type MenuItem struct {
	extension.AbstractExtension
	Spec   MenuItemSpec   `json:"spec"`
	Status MenuItemStatus `json:"status,omitempty"`
}

func (mi *MenuItem) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "",
		Version:  "v1alpha1",
		Kind:     "MenuItem",
		Plural:   "menuitems",
		Singular: "menuitem",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&MenuItem{}).GetGVK(), &MenuItem{})
}
