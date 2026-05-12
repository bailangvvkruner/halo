package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type SettingSpec struct {
	Forms []SettingForm `json:"forms"`
}

type SettingForm struct {
	FormKind string       `json:"formKind"`
	Label    string       `json:"label"`
	Items    []FormItem  `json:"items"`
}

type FormItem struct {
	Name         string      `json:"name"`
	Label        string      `json:"label"`
	Description  string      `json:"description,omitempty"`
	Type         string      `json:"type"`
	Value        interface{} `json:"value,omitempty"`
	Required     bool        `json:"required,omitempty"`
	Rules        interface{} `json:"rules,omitempty"`
	Props        interface{} `json:"props,omitempty"`
}

type SettingStatus struct {
	LastModifyTime time.Time `json:"lastModifyTime"`
}

type Setting struct {
	extension.AbstractExtension
	Spec   SettingSpec   `json:"spec"`
	Status SettingStatus `json:"status,omitempty"`
}

func (s *Setting) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "",
		Version:  "v1alpha1",
		Kind:     "Setting",
		Plural:   "settings",
		Singular: "setting",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Setting{}).GetGVK(), &Setting{})
}
