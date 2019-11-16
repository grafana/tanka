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
func walkJSON(rawDeep interface{}, path string) ([]map[string]interface{}, error) {
	flat := make([]map[string]interface{}, 0)

	// array: walkJSON for each
	if d, ok := rawDeep.([]map[string]interface{}); ok {
		for i, j := range d {
			out, err := walkJSON(j, fmt.Sprintf("%s[%v]", path, i))
			if err != nil {
				return nil, err
			}
			flat = append(flat, out...)
		}
		return flat, nil
	}

	// assert for map[string]interface{} (also aliased objx.Map)
	if m, ok := rawDeep.(objx.Map); ok {
		rawDeep = map[string]interface{}(m)
	}
	deep, ok := rawDeep.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("deep has unexpected type %T @ %s", deep, path)
	}

	// already flat?
	r := objx.New(deep)
	if r.Has("apiVersion") && r.Has("kind") {
		return []map[string]interface{}{deep}, nil
	}

	// walk it
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
