package tanka

import (
	"io/ioutil"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
)

// ListOpts are optional arguments for ListEnvs
type ListOpts struct {
	Selector labels.Selector
}

// ListEnvs returns metadata of all environments recursively found in 'dir'.
// Each directory is tested and included if it is a valid environment, either
// static or inline. If a directory is a valid environment, its subdirectories
// are not checked.
func ListEnvs(dir string, opts ListOpts) ([]*v1alpha1.Environment, error) {
	// list all environments at dir
	envs, err := list(dir)
	if err != nil {
		return nil, err
	}

	// optionally filter
	if opts.Selector == nil || opts.Selector.Empty() {
		return envs, nil
	}

	filtered := make([]*v1alpha1.Environment, 0, len(envs))
	for _, e := range envs {
		if !opts.Selector.Matches(e.Metadata) {
			continue
		}
		filtered = append(filtered, e)
	}

	return filtered, nil
}

// list implements the actual functionality described at 'ListEnvs'
func list(dir string) ([]*v1alpha1.Environment, error) {
	// list directory, also checks if dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// try if this is an env
	env, err := Peek(dir, jsonnet.Opts{})
	if err == nil {
		// it is one. don't search deeper
		return []*v1alpha1.Environment{env}, nil
	}

	// it's not one. Maybe subdirectories are?
	ch := make(chan listOut)
	routines := 0

	// recursively list in parallel
	for _, fi := range files {
		if !fi.IsDir() {
			continue
		}

		routines++
		go listShim(filepath.Join(dir, fi.Name()), ch)
	}

	// collect parallel results
	var lastErr error
	var envs []*v1alpha1.Environment

	for i := 0; i < routines; i++ {
		out := <-ch
		if out.err != nil {
			lastErr = out.err
		}

		envs = append(envs, out.envs...)
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return envs, nil
}

type listOut struct {
	envs []*v1alpha1.Environment
	err  error
}

func listShim(dir string, ch chan listOut) {
	envs, err := list(dir)
	ch <- listOut{envs: envs, err: err}
}
