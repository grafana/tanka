package tanka

import (
	"fmt"
	"sync"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

const parallel = 8

// ParseParallel evaluates multiple environments in parallel
func ParseParallel(paths []string, opts ParseParallelOpts) (envs []*v1alpha1.Environment, err error) {
	wg := sync.WaitGroup{}
	envsChan := make(chan parseJob)
	var allErrors []error

	numParallel := parallel
	if opts.Parallel > 0 {
		numParallel = opts.Parallel
	}
	for i := 0; i < numParallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs := parseWorker(envsChan)
			if errs != nil {
				allErrors = append(allErrors, errs...)
			}
		}()
	}

	results := make([]*v1alpha1.Environment, 0, len(paths))

	for _, path := range paths {
		env := &v1alpha1.Environment{}
		results = append(results, env)
		envsChan <- parseJob{
			path: path,
			opts: opts.JsonnetOpts,
			env:  env,
		}
	}
	close(envsChan)
	wg.Wait()

	for _, env := range results {
		if env == nil {
			continue
		}
		if opts.Selector == nil || opts.Selector.Empty() || opts.Selector.Matches(env.Metadata) {
			envs = append(envs, env)
		}
	}

	if len(allErrors) != 0 {
		return envs, ErrParseParallel{errors: allErrors}
	}

	return envs, nil
}

type parseJob struct {
	path string
	opts JsonnetOpts
	env  *v1alpha1.Environment
}

func parseWorker(envsChan <-chan parseJob) (errs []error) {
	for req := range envsChan {
		env, err := LoadEnvironment(req.path, Opts{JsonnetOpts: req.opts})
		if err != nil {
			errs = append(errs, fmt.Errorf("%s:\n %w", req.path, err))
			continue
		}
		*req.env = *env
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}
