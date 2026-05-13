package model

import (
	"github.com/halo-dev/halo-go/internal/extension"
)

type SystemSettings struct {
	extension.AbstractExtension
	Spec   SystemSpec  `json:"spec"`
	Status interface{} `json:"status"`
}

func (s *SystemSettings) GetGVK() extension.GVK {
	return extension.GVK{
		Group:    "",
		Version:  "v1alpha1",
		Kind:     "SystemSettings",
		Plural:   "systemsettings",
		Singular: "systemsettings",
	}
}

func init() {
	_ = extension.DefaultScheme().Register((&SystemSettings{}).GetGVK(), &SystemSettings{})
}

type SystemSpec struct {
	Basic         BasicSetting `json:"basic"`
	SEO           SEOSetting   `json:"seo"`
	CodeInjection CodeSetting   `json:"codeInjection"`
}

type BasicSetting struct {
	Title            string `json:"title"`
	Subtitle         string `json:"subtitle"`
	Logo             string `json:"logo"`
	Favicon          string `json:"favicon"`
	Language         string `json:"language"`
	TimeZone         string `json:"timeZone"`
	ExternalUrl      string `json:"externalUrl"`
	NavMode          string `json:"navMode"`
	NavItems         string `json:"navItems"`
	PostPreview      string `json:"postPreview"`
	PostListPageSize int    `json:"postListPageSize"`
	CommentEnabled   bool   `json:"commentEnabled"`
}

type SEOSetting struct {
	BlockSpiders           bool   `json:"blockSpiders"`
	Keywords               string `json:"keywords"`
	Description            string `json:"description"`
	GoogleSiteVerification string `json:"googleSiteVerification"`
	BingSiteVerification   string `json:"bingSiteVerification"`
	BaiduSiteVerification  string `json:"baiduSiteVerification"`
	QihuSiteVerification   string `json:"qihuSiteVerification"`
}

type CodeSetting struct {
	GlobalHead    string `json:"globalHead"`
	GlobalFooter  string `json:"globalFooter"`
	ContentHead   string `json:"contentHead"`
	ContentFooter string `json:"contentFooter"`
}
