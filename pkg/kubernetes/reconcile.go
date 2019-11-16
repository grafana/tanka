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

// Reconcile extracts all valid Kubernetes objects from the raw output of the
// Jsonnet compiler. A valid object is identified by the presence of `kind` and
// `apiVersion`.
// TODO: Check on `metadata.name` as well and assert that they are
// not only set but also strings.
func Reconcile(raw map[string]interface{}, spec v1alpha1.Spec, targets []*regexp.Regexp) (state manifest.List, err error) {
	docs, err := walkJSON(raw, "")
	if err != nil {
		return nil, errors.Wrap(err, "flattening manifests")
	}

	out := make(manifest.List, 0, len(docs))
	for _, d := range docs {
		o := objx.New(d)

		// complete missing namespace from spec.json
		if spec.Namespace != "" && !o.Has("metadata.namespace") {
			o.Set("metadata.namespace", spec.Namespace)
		}

		m, err := manifest.New(o)
		if err != nil {
			return nil, err
		}
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

// walkJSON traverses deeply nested kubernetes manifest and extracts them into a flat []dict.
func walkJSON(deep map[string]interface{}, path string) ([]map[string]interface{}, error) {
	r := objx.New(deep)
	if r.Has("apiVersion") && r.Has("kind") {
		return []map[string]interface{}{deep}, nil
	}

	flat := make([]map[string]interface{}, 0)

	for n, d := range deep {
		if n == "__ksonnet" {
			continue
		}
		if _, ok := d.(map[string]interface{}); !ok {
			return nil, ErrorPrimitiveReached{path, n, d}
		}
		m := objx.New(d)
		if m.Has("apiVersion") && m.Has("kind") {
			flat = append(flat, m)
		} else {
			f, err := walkJSON(m, path+"."+n)
			if err != nil {
				return nil, err
			}
			flat = append(flat, f...)
		}
	}
	return flat, nil
}

// ErrorPrimitiveReached occurs when walkJSON reaches the end of nested dicts without finding a valid Kubernetes manifest
type ErrorPrimitiveReached struct {
	path, key string
	primitive interface{}
}

func (e ErrorPrimitiveReached) Error() string {
	return fmt.Sprintf("recursion did not resolve in a valid Kubernetes object, "+
		"because one of `kind` or `apiVersion` is missing in path `.%s`."+
		" Found non-dict value `%s` of type `%T` instead.",
		e.path, e.key, e.primitive)
}
