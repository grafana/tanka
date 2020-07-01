package process

import (
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

const (
	// AnnotationNamespaced can be set on any resource to override the decision
	// whether 'metadata.namespace' is set by Tanka
	AnnotationNamespaced = MetadataPrefix + "/namespaced"
)

// Namespace injects the default namespace of the environment into each
// resources, that does not already define one. AnnotationNamespaced can be used
// to disable this per resource
func Namespace(list manifest.List, def string) manifest.List {
	if def == "" {
		return list
	}

	for i, m := range list {
		namespaced := true

		// check for annotation override
		if s, ok := m.Metadata().Annotations()[AnnotationNamespaced]; ok {
			namespaced = s == "true"
		}

		if namespaced && !m.Metadata().HasNamespace() {
			m.Metadata()["namespace"] = def
		}

		// remove annotations if empty (we always create those by accessing them)
		if len(m.Metadata().Annotations()) == 0 {
			delete(m.Metadata(), "annotations")
		}

		list[i] = m
	}

	return list
}
