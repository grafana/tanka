package tanka

import (
	"io/ioutil"
	"path/filepath"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

type ListOpts struct {
	Selector labels.Selector
}

func ListEnvs(dir string, opts ListOpts) ([]*v1alpha1.Environment, error) {
	data, err := envs(dir)
	if err != nil {
		return nil, err
	}

	out := make([]*v1alpha1.Environment, 0, len(data))
	for _, e := range data {
		if opts.Selector == nil || opts.Selector.Empty() || opts.Selector.Matches(e.Metadata) {
			out = append(out, e)
		}
	}

	return out, nil
}

func envs(dir string) ([]*v1alpha1.Environment, error) {
	// check if dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// try if this is an env
	env, err := Peek(dir, jsonnet.Opts{})
	if err == nil {
		return []*v1alpha1.Environment{env}, nil
	}

	// it's not one. Maybe it has children?
	var out []*v1alpha1.Environment

	errCh := make(chan error)
	envCh := make(chan []*v1alpha1.Environment)
	routines := 0

	for _, fi := range files {
		if !fi.IsDir() {
			continue
		}

		routines++
		go envP(filepath.Join(dir, fi.Name()), envCh, errCh)
	}

	var lastErr error

	for i := 0; i < routines; i++ {
		select {
		case envs := <-envCh:
			out = append(out, envs...)
		case err := <-errCh:
			lastErr = err
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return out, nil
}

func envP(dir string, c chan []*v1alpha1.Environment, e chan error) {
	data, err := envs(dir)
	if err != nil {
		e <- err
	}
	c <- data
}
