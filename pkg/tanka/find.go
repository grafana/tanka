package tanka

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
)

// FindOpts are optional arguments for FindEnvs
type FindOpts struct {
	JsonnetOpts
	Selector labels.Selector
}

// FindEnvs returns metadata of all environments recursively found in 'path'.
// Each directory is tested and included if it is a valid environment, either
// static or inline. If a directory is a valid environment, its subdirectories
// are not checked.
func FindEnvs(path string, opts FindOpts) ([]*v1alpha1.Environment, error) {
	// find all environments at dir
	envs, err := find(path, Opts{JsonnetOpts: opts.JsonnetOpts})
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

// find implements the actual functionality described at 'FindEnvs'
func find(path string, opts Opts) ([]*v1alpha1.Environment, error) {
	// try if this has envs
	list, err := List(path, opts)
	if len(list) != 0 && err == nil {
		// it has. don't search deeper
		return list, nil
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// if path is a file, don't search deeper
	if !stat.IsDir() {
		return nil, nil
	}

	// list directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// it's not one. Maybe subdirectories are?
	ch := make(chan findOut)
	routines := 0

	// recursively find in parallel
	for _, fi := range files {
		if !fi.IsDir() {
			continue
		}

		routines++
		go findShim(filepath.Join(path, fi.Name()), opts, ch)
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

type findOut struct {
	envs []*v1alpha1.Environment
	err  error
}

func findShim(dir string, opts Opts, ch chan findOut) {
	envs, err := find(dir, opts)
	ch <- findOut{envs: envs, err: err}
}
