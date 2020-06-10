package process

import (
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/resources"
)

const (
	// AnnotationNamespaced can be set on any resource to override the decision
	// whether 'metadata.namespace' is set by Tanka
	AnnotationNamespaced = MetadataPrefix + "/namespaced"
)

// Namespace adds `metadata.namespace` fields to namespaced resources. If a
// resource resource is namespaced is discovered based on data from `kubectl
// api-resources`, or from the tanka.dev/namespaced annotation if present.
func Namespace(list manifest.List, ns string, r resources.Store) manifest.List {
	if ns == "" {
		return list
	}

	for i, m := range list {
		namespaced := false

		// check for annotation override
		if s, ok := m.Metadata().Annotations()[AnnotationNamespaced]; ok {
			namespaced = s == "true"
		} else {
			namespaced = r.Namespaced(m)
		}

		if namespaced && !m.Metadata().HasNamespace() {
			m.Metadata()["namespace"] = ns
		}

		// remove annotations if empty (we always create those by accessing them)
		if len(m.Metadata().Annotations()) == 0 {
			delete(m.Metadata(), "annotations")
		}

		list[i] = m
	}

	return list
}
