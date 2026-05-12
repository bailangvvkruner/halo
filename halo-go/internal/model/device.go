package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type DeviceSpec struct {
	Platform string `json:"platform"`
	Os       string `json:"os"`
	Vm       string `json:"vm"`
	Browser  string `json:"browser"`
	IP       string `json:"ip"`
}

type DeviceStatus struct {
	LastAccessedTime time.Time `json:"lastAccessedTime,omitempty"`
}

type Device struct {
	extension.AbstractExtension `json:",inline"`
	Spec                        DeviceSpec   `json:"spec"`
	Status                      DeviceStatus `json:"status,omitempty"`
}

func newDevice() *Device { return &Device{} }

func (d *Device) GetGVK() extension.GVK {
	return extension.NewGVK("", "v1alpha1", "Device", "devices", "device")
}

func init() {
	_ = extension.DefaultScheme().Register((&Device{}).GetGVK(), &Device{})
}
