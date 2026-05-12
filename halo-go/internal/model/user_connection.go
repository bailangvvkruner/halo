package model

import (
	"github.com/halo-dev/halo-go/internal/extension"
)

type UserConnectionSpec struct {
	Username           string            `json:"username"`
	Provider           string            `json:"provider"`
	ProviderId         string            `json:"providerId"`
	ProviderUsername   string            `json:"providerUsername,omitempty"`
	ProviderDisplayName string           `json:"providerDisplayName,omitempty"`
	UserData           map[string]string `json:"userData,omitempty"`
}

type UserConnection struct {
	extension.AbstractExtension `json:",inline"`
	Spec                        UserConnectionSpec `json:"spec"`
}

func newUserConnection() *UserConnection { return &UserConnection{} }

func (uc *UserConnection) GetGVK() extension.GVK {
	return extension.NewGVK("", "v1alpha1", "UserConnection", "userconnections", "userconnection")
}
