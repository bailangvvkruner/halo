package model

import (
	"github.com/halo-dev/halo-go/internal/extension"
)

type AnnotationDefinition struct {
	Key           string            `json:"key"`
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	DefaultValues map[string]string `json:"defaultValues,omitempty"`
	FormSchema    []any             `json:"formSchema,omitempty"`
}

type AnnotationSettingSpec struct {
	Definitions []AnnotationDefinition `json:"definitions,omitempty"`
}

type AnnotationSetting struct {
	extension.AbstractExtension `json:",inline"`
	Spec                        AnnotationSettingSpec `json:"spec"`
}

func newAnnotationSetting() *AnnotationSetting { return &AnnotationSetting{} }

func (as *AnnotationSetting) GetGVK() extension.GVK {
	return extension.NewGVK("", "v1alpha1", "AnnotationSetting", "annotationsettings", "annotationsetting")
}

func init() {
	_ = extension.DefaultScheme().Register((&AnnotationSetting{}).GetGVK(), &AnnotationSetting{})
}
