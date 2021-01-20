package tanka

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

const defaultParallelism = 8

type parallelOpts struct {
	Opts
	Selector    labels.Selector
	Parallelism int
}

// parallelLoadEnvironments evaluates multiple environments in parallel
func parallelLoadEnvironments(paths []string, opts parallelOpts) ([]*v1alpha1.Environment, error) {
	jobsCh := make(chan parallelJob)
	list := make(map[string]string)
	for _, path := range paths {
		envs, err := FindEnvs(path, FindOpts{opts.Selector})
		if err != nil {
			return nil, err
		}
		for _, env := range envs {
			if opts.Name != "" && opts.Name != env.Metadata.Name {
				continue
			}
			list[env.Metadata.Name] = path
		}
	}
	outCh := make(chan parallelOut, len(list))

	if opts.Parallelism <= 0 {
		opts.Parallelism = defaultParallelism
	}

	for i := 0; i < opts.Parallelism; i++ {
		go parallelWorker(jobsCh, outCh)
	}

	for name, path := range list {
		o := opts.Opts

		// TODO: This is required because the map[string]string in here is not
		// concurrency-safe. Instead of putting this burden on the caller, find
		// a way to handle this inside the jsonnet package. A possible way would
		// be to make the jsonnet package less general, more tightly coupling it
		// to Tanka workflow thus being able to handle such cases
		o.JsonnetOpts = o.JsonnetOpts.Clone()

		o.Name = name
		jobsCh <- parallelJob{
			path: path,
			opts: o,
		}
	}
	close(jobsCh)

	var envs []*v1alpha1.Environment
	var errors []error
	for i := 0; i < len(list); i++ {
		out := <-outCh
		if out.err != nil {
			errors = append(errors, out.err)
			continue
		}
		if opts.Selector == nil || opts.Selector.Empty() || opts.Selector.Matches(out.env.Metadata) {
			envs = append(envs, out.env)
		}
	}

	if len(errors) != 0 {
		return envs, ErrParallel{errors: errors}
	}

	return envs, nil
}

type parallelJob struct {
	path string
	opts Opts
}

type parallelOut struct {
	env *v1alpha1.Environment
	err error
}

func parallelWorker(jobsCh <-chan parallelJob, outCh chan parallelOut) {
	for job := range jobsCh {
		env, err := LoadEnvironment(job.path, job.opts)
		if err != nil {
			err = fmt.Errorf("%s:\n %w", job.path, err)
		}
		outCh <- parallelOut{env: env, err: err}
	}
}
