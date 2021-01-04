package tanka

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

const baseDirIndicator = "main.jsonnet"
const parallel = 8

// FindBaseDirs searches for possible environments
func FindBaseDirs(path string) (paths []string, err error) {
	pathInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !pathInfo.IsDir() {
		return append(paths, path), nil
	}

	_, _, _, err = jpath.Resolve(path)
	if err == jpath.ErrorNoRoot {
		return nil, err
	}

	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(filepath.Join(path, baseDirIndicator)); err != nil {
			// missing file, not a valid environment directory
			return nil
		}
		paths = append(paths, path)
		return nil
	}); err != nil {
		return nil, err
	}
	return paths, nil
}

// FindEnvironments searches for actual environments
// ignores directories if no environments are found
func FindEnvironments(path string, selector labels.Selector) (envs map[string][]*v1alpha1.Environment, err error) {
	opts := ParallelOpts{
		JsonnetOpts: JsonnetOpts{
			EvalScript: EnvsOnlyEvalScript,
		},
		Selector: selector,
	}

	paths, err := FindBaseDirs(path)
	if err != nil {
		return nil, err
	}

	envs, err = LoadEnvironmentsParallel(paths, opts)

	if err != nil {
		switch err.(type) {
		case ErrParseParallel:
			// ignore ErrNoEnv errors
			e := err.(ErrParseParallel)
			var errors []error
			for _, err := range e.Errors {
				switch err.(type) {
				case ErrNoEnv:
					continue
				default:
					errors = append(errors, err)
				}
			}
			if len(errors) != 0 {
				return nil, ErrParseParallel{Errors: errors}
			}
		default:
			return nil, err
		}
	}

	return envs, nil
}

// LoadEnvironmentsParallel evaluates multiple environments in parallel
func LoadEnvironmentsParallel(paths []string, opts ParallelOpts) (envs map[string][]*v1alpha1.Environment, err error) {
	envs = make(map[string][]*v1alpha1.Environment, 0)

	wg := sync.WaitGroup{}
	envsChan := make(chan loadJob)
	var allErrors []error

	numParallel := parallel
	if opts.Parallel > 0 {
		numParallel = opts.Parallel
	}
	for i := 0; i < numParallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, errs := loadWorker(envsChan)
			if errs != nil {
				allErrors = append(allErrors, errs...)
			}
			if result != nil {
				for path, v := range result {
					for _, env := range v {
						if env == nil {
							continue
						}
						if opts.Selector == nil || opts.Selector.Empty() || opts.Selector.Matches(env.Metadata) {
							envs[path] = append(envs[path], env)
						}
					}
				}
			}
		}()
	}

	for _, path := range paths {
		envsChan <- loadJob{
			path: path,
			opts: opts.JsonnetOpts,
		}
	}
	close(envsChan)
	wg.Wait()

	if len(allErrors) != 0 {
		return envs, ErrParseParallel{Errors: allErrors}
	}

	return envs, nil
}

type loadJob struct {
	path string
	opts JsonnetOpts
	env  *v1alpha1.Environment
}

func loadWorker(envsChan <-chan loadJob) (envs map[string][]*v1alpha1.Environment, errs []error) {
	envs = make(map[string][]*v1alpha1.Environment, 0)
	for req := range envsChan {
		loaded, err := LoadEnvironments(req.path, Opts{JsonnetOpts: req.opts})
		if err != nil {
			errs = append(errs, fmt.Errorf("%s:\n %w", req.path, err))
			continue
		}
		if envs[req.path] == nil {
			envs[req.path] = make([]*v1alpha1.Environment, 0)
		}
		envs[req.path] = append(envs[req.path], loaded...)
	}
	if len(errs) != 0 {
		return envs, errs
	}
	return envs, nil
}
