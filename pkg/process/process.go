package process

import (
	"errors"
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

const (
	MetadataPrefix   = "tanka.dev"
	LabelEnvironment = MetadataPrefix + "/environment"
)

// Process converts the raw Jsonnet evaluation result (JSON tree) into a flat
// list of Kubernetes objects, also applying some transformations:
// - tanka.dev/** labels
// - filtering
// - best-effort sorting
func Process(raw map[string]interface{}, cfg v1alpha1.Config, exprs Matchers) (manifest.List, error) {
	// Scan for everything that looks like a Kubernetes object
	extracted, err := Extract(raw)
	if err != nil {
		return nil, err
	}

	// Unwrap *List types
	if err := Unwrap(extracted); err != nil {
		return nil, err
	}

	out := make(manifest.List, 0, len(extracted))
	for _, m := range extracted {
		out = append(out, m)
	}

	// tanka.dev/** labels
	out = Label(out, cfg)

	// Perhaps filter for kind/name expressions
	if len(exprs) > 0 {
		out = Filter(out, exprs)
	}

	// Best-effort dependency sort
	Sort(out)

	return out, nil
}

// Label conditionally adds tanka.dev/** labels to each manifest in the List
func Label(list manifest.List, cfg v1alpha1.Config) manifest.List {
	for i, m := range list {
		// inject tanka.dev/environment label
		if cfg.Spec.InjectLabels {
			m.Metadata().Labels()[LabelEnvironment] = cfg.Metadata.NameLabel()
		}
		list[i] = m
	}

	return list
}

// Unwrap returns all Kubernetes objects in the manifest. If m is not a List
// type, a one item List is returned
func Unwrap(manifests map[string]manifest.Manifest) error {
	for path, m := range manifests {
		if !m.IsList() {
			continue
		}

		items, err := m.Items()
		if err != nil {
			return err
		}

		for index, i := range items {
			name := fmt.Sprintf("%s.items[%v]", path, index)

			var e *manifest.SchemaError
			if errors.As(i.Verify(), &e) {
				e.Name = name
				return e
			}

			manifests[name] = i
		}

		delete(manifests, path)
	}

	return nil
}
