package tanka

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/pkg/errors"
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
	envs, errs := find(path, Opts{JsonnetOpts: opts.JsonnetOpts})
	if errs != nil {
		return envs, ErrParallel{errors: errs}
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

func findErr(path string, err error) []error {
	return []error{fmt.Errorf("%s:\n %w", path, err)}
}

// find implements the actual functionality described at 'FindEnvs'
func find(path string, opts Opts) ([]*v1alpha1.Environment, []error) {
	// try if this has envs
	list, err := List(path, opts)
	if err != nil &&
		// expected when looking for environments
		!errors.As(err, &jpath.ErrorNoBase{}) &&
		!errors.As(err, &jpath.ErrorFileNotFound{}) {
		return nil, findErr(path, err)
	}
	if len(list) != 0 {
		// it has. don't search deeper
		return list, nil
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, findErr(path, err)
	}

	// if path is a file, don't search deeper
	if !stat.IsDir() {
		return nil, nil
	}

	// list directory
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, findErr(path, err)
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
	var errs []error
	var envs []*v1alpha1.Environment

	for i := 0; i < routines; i++ {
		out := <-ch
		if out.errs != nil {
			errs = append(errs, out.errs...)
		}

		envs = append(envs, out.envs...)
	}

	if len(errs) != 0 {
		return envs, errs
	}

	return envs, nil
}

type findOut struct {
	envs []*v1alpha1.Environment
	errs []error
}

func findShim(dir string, opts Opts, ch chan findOut) {
	envs, errs := find(dir, opts)
	ch <- findOut{envs: envs, errs: errs}
}
