package tanka

import (
	"fmt"
	"log"
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
func FindBaseDirs(workdir string) (dirs []string, err error) {
	_, _, _, err = jpath.Resolve(workdir)
	if err == jpath.ErrorNoRoot {
		return nil, err
	}

	if err := filepath.Walk(workdir, func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(filepath.Join(path, baseDirIndicator)); err != nil {
			// missing file, not a valid environment directory
			return nil
		}
		dirs = append(dirs, path)
		return nil
	}); err != nil {
		return nil, err
	}
	return dirs, nil
}

// FindEnvironments searches for actual environments
// ignores directories if no environments are found
func FindEnvironments(workdir string, selector labels.Selector) (envs []*v1alpha1.Environment, err error) {
	dirs, err := FindBaseDirs(workdir)
	if err != nil {
		return nil, err
	}
	opts := ParseParallelOpts{
		JsonnetOpts: JsonnetOpts{
			EvalScript: EnvsOnlyEvalScript,
		},
		Selector: selector,
	}
	envs, err = ParseParallel(dirs, opts)

	if err != nil {
		switch err.(type) {
		case ErrParseParallel:
			// ignore ErrNoEnv errors
			e := err.(ErrParseParallel)
			var errors []error
			for _, err := range e.errors {
				switch err.(type) {
				case ErrNoEnv:
					continue
				default:
					errors = append(errors, err)
				}
			}
			if len(errors) != 0 {
				return nil, ErrParseParallel{errors: errors}
			}
		default:
			return nil, err
		}
	}

	return envs, nil
}

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
		log.Printf("Parsing %s\n", req.path)
		_, env, err := ParseEnv(req.path, req.opts)
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
