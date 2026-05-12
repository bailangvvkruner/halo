package extension

import (
	"fmt"
	"strings"
)

type GVK struct {
	Group    string `json:"group"`
	Version  string `json:"version"`
	Kind     string `json:"kind"`
	Plural   string `json:"plural"`
	Singular string `json:"singular"`
}

func NewGVK(group, version, kind, plural, singular string) GVK {
	return GVK{
		Group:    group,
		Version:  version,
		Kind:     kind,
		Plural:   plural,
		Singular: singular,
	}
}

func (g GVK) String() string {
	return fmt.Sprintf("%s/%s/%s", g.Group, g.Version, g.Kind)
}

func (g GVK) SortKey() string {
	return strings.ToLower(g.String())
}

func (g GVK) Equal(other GVK) bool {
	return g.Group == other.Group && g.Version == other.Version && g.Kind == other.Kind
}

func (g GVK) IsEmpty() bool {
	return g.Group == "" && g.Version == "" && g.Kind == ""
}

func (g GVK) ResourceIdentifier() string {
	return fmt.Sprintf("%s/%s/%s", g.Group, g.Version, g.Plural)
}
