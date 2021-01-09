package tanka

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

const defaultParallelism = 8

type parallelOpts struct {
	JsonnetOpts
	Selector    labels.Selector
	Parallelism int
}

// parallelLoadEnvironments evaluates multiple environments in parallel
func parallelLoadEnvironments(paths []string, opts parallelOpts) ([]*v1alpha1.Environment, error) {
	wg := sync.WaitGroup{}
	jobsCh := make(chan parallelJob)

	if opts.Parallelism <= 0 {
		opts.Parallelism = defaultParallelism
	}

	for i := 0; i < opts.Parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parallelWorker(jobsCh)
		}()
	}

	var results []*parallelOut
	for _, path := range paths {
		out := &parallelOut{}
		results = append(results, out)
		jobsCh <- parallelJob{
			path: path,
			opts: Opts{JsonnetOpts: opts.JsonnetOpts},
			out:  out,
		}
	}
	close(jobsCh)

	var envs []*v1alpha1.Environment
	var errors []error
	for _, out := range results {
		if out.err != nil {
			errors = append(errors, out.err)
			continue
		}
		if opts.Selector == nil || opts.Selector.Empty() || opts.Selector.Matches(out.env.Metadata) {
			envs = append(envs, &out.env)
		}
	}
	wg.Wait()

	if len(errors) != 0 {
		return envs, ErrParallel{errors: errors}
	}

	return envs, nil
}

type parallelJob struct {
	path string
	opts Opts
	out  *parallelOut
}

type parallelOut struct {
	env v1alpha1.Environment
	err error
}

func parallelWorker(jobsCh <-chan parallelJob) {
	for job := range jobsCh {
		env, err := LoadEnvironment(job.path, job.opts)
		if err != nil {
			err = fmt.Errorf("%s:\n %w", job.path, err)
		}
		*job.out = parallelOut{*env, err}
	}
}
