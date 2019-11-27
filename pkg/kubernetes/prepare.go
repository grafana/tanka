package kubernetes

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	funk "github.com/thoas/go-funk"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Labels injected into the manifest
const (
	// LabelPrefix is a string to mark all of Tanka's fields as such
	LabelPrefix = "tanka.dev"

	// LabelJSONPath is location where the manifest is found in the jsonnet output
	LabelJSONPath = LabelPrefix + "/jsonpath"
)

// Prepare extracts all valid Kubernetes objects from the raw output of the
// Jsonnet compiler. A valid object is identified by the presence of `kind` and
// `apiVersion`.
// TODO: Check on `metadata.name` as well and assert that they are
// not only set but also strings.
// It also injects Tanka specific labels
func Prepare(raw map[string]interface{}, spec v1alpha1.Spec, targets []*regexp.Regexp) (state manifest.List, err error) {
	// extract from jsonnet output
	extracted, err := extract(raw)
	if err != nil {
		return nil, errors.Wrap(err, "flattening manifests")
	}

	out := make(manifest.List, 0, len(extracted))
	for jsonpath, m := range extracted {
		// default namespace if unset
		if spec.Namespace != "" && !m.Metadata().HasNamespace() {
			m.Metadata()["namespace"] = spec.Namespace
		}

		// add jsonpath to annotations
		m.Metadata().Annotations()[LabelJSONPath] = jsonpath

		out = append(out, m)
	}

	// optionally filter the working set of objects
	if len(targets) > 0 {
		tmp := funk.Filter(out, func(i interface{}) bool {
			p := objectspec(i.(manifest.Manifest))
			for _, t := range targets {
				if t.MatchString(strings.ToLower(p)) {
					return true
				}
			}
			return false
		}).([]manifest.Manifest)
		out = manifest.List(tmp)
	}

	// Stable output order
	sort.SliceStable(out, func(i int, j int) bool {
		if out[i].Kind() != out[j].Kind() {
			return out[i].Kind() < out[j].Kind()
		}
		return out[i].Metadata().Name() < out[j].Metadata().Name()
	})

	return out, nil
}

func extract(deep interface{}) (map[string]manifest.Manifest, error) {
	extracted := make(map[string]manifest.Manifest)
	if err := walkJSON(deep, extracted, nil); err != nil {
		return nil, err
	}
	return extracted, nil
}

// walkJSON traverses deeply nested kubernetes manifest and extracts them into a flat []dict.
func walkJSON(deep interface{}, extracted map[string]manifest.Manifest, path trace) error {
	// array: walkJSON for each
	if d, ok := deep.([]map[string]interface{}); ok {
		for i, j := range d {
			path := append(path, fmt.Sprintf("[%v]", i))
			if err := walkJSON(j, extracted, path); err != nil {
				return err
			}
		}
		return nil
	}

	// assert for map[string]interface{} (also aliased objx.Map)
	if m, ok := deep.(objx.Map); ok {
		deep = map[string]interface{}(m)
	}
	deep, ok := deep.(map[string]interface{})
	if !ok {
		return fmt.Errorf("deep has unexpected type %T @ %s", deep, path)
	}

	// already flat?
	r := objx.New(deep)

	if r.Has("apiVersion") && r.Has("kind") {
		extracted[path.Full()] = deep.(map[string]interface{})
		return nil
	}

	// walk it
	for key, d := range deep.(map[string]interface{}) {
		if key == "__ksonnet" {
			continue
		}
		path := append(path, key)

		if _, ok := d.(map[string]interface{}); !ok {
			return ErrorPrimitiveReached{path.Base(), key, d}
		}

		m := objx.New(d)
		if m.Has("apiVersion") && m.Has("kind") {
			mf, err := manifest.NewFromObj(m)
			if err != nil {
				return err.(*manifest.SchemaError).WithName(path.Full())
			}
			extracted[path.Full()] = mf
		} else {
			if err := walkJSON(m, extracted, path); err != nil {
				return err
			}
		}
	}
	return nil
}

type trace []string

func (t trace) Full() string {
	return "." + strings.Join(t, ".")
}

func (t trace) Base() string {
	if len(t) > 0 {
		t = t[:len(t)-1]
	}
	return "." + strings.Join(t, ".")
}

// ErrorPrimitiveReached occurs when walkJSON reaches the end of nested dicts without finding a valid Kubernetes manifest
type ErrorPrimitiveReached struct {
	path, key string
	primitive interface{}
}

func (e ErrorPrimitiveReached) Error() string {
	return fmt.Sprintf("recursion did not resolve in a valid Kubernetes object, "+
		"because one of `kind` or `apiVersion` is missing in path `%s`."+
		" Found non-dict value `%s` of type `%T` instead.",
		e.path, e.key, e.primitive)
}
