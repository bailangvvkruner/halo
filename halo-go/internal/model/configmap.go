package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type ConfigMapSpec struct {
	Data map[string]string `json:"data"`
}

type ConfigMapStatus struct {
	LastModifyTime time.Time `json:"lastModifyTime"`
}

type ConfigMap struct {
	extension.AbstractExtension
	Spec   ConfigMapSpec   `json:"spec"`
	Status ConfigMapStatus `json:"status,omitempty"`
}

func (cm *ConfigMap) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "",
		Version:  "v1alpha1",
		Kind:     "ConfigMap",
		Plural:   "configmaps",
		Singular: "configmap",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&ConfigMap{}).GetGVK(), &ConfigMap{})
}
