package resources

//go:generate go run ./gen

import (
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// StaticStore is a pre-compiled list of API Resources based on the most recent
// Kubernetes version in a default configuration

// Store holds information about API resources known to a Kubernetes cluster
type Store []Resource

// Namespaced returns whether a resource is namespace-specific or cluster-wide
func (s Store) Namespaced(m manifest.Manifest) bool {
	for _, res := range s {
		if m.Kind() == res.Kind {
			return res.Namespaced
		}
	}

	return false
}

// Resource is a Kubernetes API Resource
type Resource struct {
	APIGroup   string `json:"APIGROUP"`
	Kind       string `json:"KIND"`
	Name       string `json:"NAME"`
	Namespaced bool   `json:"NAMESPACED,string"`
	Shortnames string `json:"SHORTNAMES"`
	Verbs      string `json:"VERBS"`
}

func (r Resource) FQN() string {
	return strings.TrimSuffix(r.Kind+"."+r.APIGroup, ".")
}
