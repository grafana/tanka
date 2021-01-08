package tanka

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

const parallel = 8

type parallelOpts struct {
	JsonnetOpts JsonnetOpts
	Selector    labels.Selector
	Parallel    int
}

// parallelLoadEnvironments evaluates multiple environments in parallel
func parallelLoadEnvironments(paths []string, opts parallelOpts) ([]*v1alpha1.Environment, error) {
	wg := sync.WaitGroup{}
	envsChan := make(chan parallelJob)
	outChan := make(chan parallelOut)

	if opts.Parallel <= 0 {
		opts.Parallel = parallel
	}

	for i := 0; i < opts.Parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parallelWorker(envsChan, outChan)
		}()
	}

	jobs := 0
	for _, path := range paths {
		envsChan <- parallelJob{
			path: path,
			opts: Opts{JsonnetOpts: opts.JsonnetOpts},
		}
		jobs++
	}
	close(envsChan)

	var envs []*v1alpha1.Environment
	var errors []error
	for i := 0; i < jobs; i++ {
		out := <-outChan
		if out.err != nil {
			errors = append(errors, out.err)
		}
		if out.env == nil {
			continue
		}
		if opts.Selector == nil || opts.Selector.Empty() || opts.Selector.Matches(out.env.Metadata) {
			envs = append(envs, out.env)
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
}

type parallelOut struct {
	env *v1alpha1.Environment
	err error
}

func parallelWorker(jobs <-chan parallelJob, out chan parallelOut) {
	for job := range jobs {
		env, err := LoadEnvironment(job.path, job.opts)
		if err != nil {
			err = fmt.Errorf("%s:\n %w", job.path, err)
		}
		out <- parallelOut{env: env, err: err}
	}
}
