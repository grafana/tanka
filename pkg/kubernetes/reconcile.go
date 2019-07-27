package kubernetes

import (
	"fmt"

	"github.com/stretchr/objx"
)

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

// walkJSON traverses deeply nested kubernetes manifest and extracts them into a flat []dict.
func walkJSON(deep map[string]interface{}, path string) ([]Manifest, error) {
	flat := []Manifest{}

	for n, d := range deep {
		if n == "__ksonnet" {
			continue
		}
		if _, ok := d.(map[string]interface{}); !ok {
			return nil, ErrorPrimitiveReached{path, n, d}
		}
		m := objx.New(d)
		if m.Has("apiVersion") && m.Has("kind") {
			flat = append(flat, Manifest(m))
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
