package model

import "github.com/halo-dev/halo-go/internal/extension"

type ThemeSpec struct {
	DisplayName   string `json:"displayName"`
	Description   string `json:"description,omitempty"`
	Version       string `json:"version"`
	Author        string `json:"author,omitempty"`
	Logo          string `json:"logo,omitempty"`
	Homepage      string `json:"homepage,omitempty"`
	Repo          string `json:"repo,omitempty"`
	Requires      string `json:"requires,omitempty"`
	SettingName   string `json:"settingName,omitempty"`
	Active        bool   `json:"active"`
	Installed     bool   `json:"installed,omitempty"`
	ConfigMapName string `json:"configMapName,omitempty"`
}

type ThemeStatus struct {
	Phase string `json:"phase,omitempty"`
}

type Theme struct {
	extension.AbstractExtension
	Spec   ThemeSpec   `json:"spec"`
	Status ThemeStatus `json:"status,omitempty"`
}

func (t *Theme) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "theme.halo.run",
		Version:  "v1alpha1",
		Kind:     "Theme",
		Plural:   "themes",
		Singular: "theme",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&Theme{}).GetGVK(), &Theme{})
}
