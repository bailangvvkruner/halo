package model

import (
	"time"

	"github.com/halo-dev/halo-go/internal/extension"
)

type UserSpec struct {
	UserName      string     `json:"userName"`
	Email         string     `json:"email"`
	DisplayName   string     `json:"displayName"`
	Avatar        string     `json:"avatar,omitempty"`
	Bio           string     `json:"bio,omitempty"`
	Password      string     `json:"-"`
	RawPassword   string     `json:"password,omitempty"`
	LastLoginTime *time.Time `json:"lastLoginTime,omitempty"`
}

type UserStatus struct {
	Permalink string `json:"permalink"`
}

type User struct {
	extension.AbstractExtension
	Spec   UserSpec   `json:"spec"`
	Status UserStatus `json:"status,omitempty"`
}

func (u *User) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "",
		Version:  "v1alpha1",
		Kind:     "User",
		Plural:   "users",
		Singular: "user",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&User{}).GetGVK(), &User{})
}
