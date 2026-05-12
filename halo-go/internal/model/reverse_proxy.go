package model

import (
	"github.com/halo-dev/halo-go/internal/extension"
)

type ReverseProxyRule struct {
	Path          string            `json:"path"`
	HeaderFilters map[string]string `json:"headerFilters,omitempty"`
	TargetUrl     string            `json:"targetUrl"`
	AuthStrategy  string            `json:"authStrategy,omitempty"`
	CachingEnabled bool             `json:"cachingEnabled"`
	CacheTTL      int               `json:"cacheTTL,omitempty"`
}

type ReverseProxySpec struct {
	Rules []ReverseProxyRule `json:"rules,omitempty"`
}

type ReverseProxy struct {
	extension.AbstractExtension `json:",inline"`
	Spec                        ReverseProxySpec `json:"spec"`
}

func newReverseProxy() *ReverseProxy { return &ReverseProxy{} }

func (rp *ReverseProxy) GetGVK() extension.GVK {
	return extension.NewGVK("", "v1alpha1", "ReverseProxy", "reverseproxies", "reverseproxy")
}
