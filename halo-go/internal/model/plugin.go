package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type PluginSpec struct {
	DisplayName  string     `json:"displayName"`
	Description  string     `json:"description,omitempty"`
	Version      string     `json:"version"`
	Author       string     `json:"author,omitempty"`
	Logo         string     `json:"logo,omitempty"`
	Homepage     string     `json:"homepage,omitempty"`
	Repo         string     `json:"repo,omitempty"`
	Requires     string     `json:"requires,omitempty"`
	SettingName  string     `json:"settingName,omitempty"`
	Enabled      bool       `json:"enabled"`
	SetupTime    *time.Time `json:"setupTime,omitempty"`
	StartTime    *time.Time `json:"startTime,omitempty"`
}

type PluginStatus struct {
	Phase string `json:"phase,omitempty"`
}

type Plugin struct {
	extension.AbstractExtension
	Spec   PluginSpec   `json:"spec"`
	Status PluginStatus `json:"status,omitempty"`
}

func (p *Plugin) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "plugin.halo.run",
		Version:  "v1alpha1",
		Kind:     "Plugin",
		Plural:   "plugins",
		Singular: "plugin",
	}
}

func init() {
	scheme := extension.NewScheme()
	_ = scheme.Register((&Plugin{}).GetGVK(), &Plugin{})
}
