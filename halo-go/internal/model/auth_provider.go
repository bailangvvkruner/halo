package model

import (
	"github.com/halo-dev/halo-go/internal/extension"
)

type AuthProviderSpec struct {
	DisplayName   string `json:"displayName"`
	Type          string `json:"type"`
	Enabled       bool   `json:"enabled"`
	ConfigMapName string `json:"configMapName,omitempty"`
}

type AuthProvider struct {
	extension.AbstractExtension `json:",inline"`
	Spec                        AuthProviderSpec `json:"spec"`
}

func newAuthProvider() *AuthProvider { return &AuthProvider{} }

func (ap *AuthProvider) GetGVK() extension.GVK {
	return extension.NewGVK("", "v1alpha1", "AuthProvider", "authproviders", "authprovider")
}
