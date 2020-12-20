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

const BASEDIR_INDICATOR = "main.jsonnet"
const PARALLEL = 8

// FindBaseDirs searches for possible environments
func FindBaseDirs(workdir string) (dirs []string, err error) {
	_, _, _, err = jpath.Resolve(workdir)
	if err == jpath.ErrorNoRoot {
		return nil, err
	}

	if err := filepath.Walk(workdir, func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(filepath.Join(path, BASEDIR_INDICATOR)); err != nil {
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
	opts := ParseOpts{
		Evaluator: EnvsOnlyEvaluator,
		Selector:  selector,
	}
	envs, errs := ParseEnvs(dirs, opts)

	var returnErrs string
	for _, err := range errs {
		switch err.(type) {
		case ErrNoEnv:
			continue
		default:
			fmt.Print(err)
			returnErrs = fmt.Sprintf("%s\n%s", returnErrs, err)
		}
	}
	if len(returnErrs) != 0 {
		return nil, fmt.Errorf("Unable to parse selected Environments: \n%s", returnErrs)
	}

	return envs, nil
}

// ParseEnvs evaluates multiple environments in parallel
func ParseEnvs(paths []string, opts ParseOpts) (envs []*v1alpha1.Environment, errs []error) {
	wg := sync.WaitGroup{}
	envsChan := make(chan parseEnvsRoutineOpts)
	var allErrors []error

	numParallel := PARALLEL
	if opts.Parallel > 0 {
		numParallel = opts.Parallel
	}
	for i := 0; i < numParallel; i++ {
		wg.Add(1)
		go func() {
			err := parseEnvsRoutine(envsChan)
			if err != nil {
				allErrors = append(allErrors, err)
			}
			wg.Done()
		}()
	}

	results := make([]*v1alpha1.Environment, len(paths))
	currentIndex := 0

	for _, path := range paths {
		env := &v1alpha1.Environment{}
		results[currentIndex] = env
		currentIndex++
		envsChan <- parseEnvsRoutineOpts{
			path: path,
			opts: opts,
			env:  env,
		}
	}
	close(envsChan)
	wg.Wait()

	for _, env := range results {
		if env != nil {
			if opts.Selector == nil || opts.Selector.Empty() || opts.Selector.Matches(env.Metadata) {
				envs = append(envs, env)
			}
		}
	}

	if len(allErrors) != 0 {
		return envs, allErrors
	}

	return envs, nil
}

type parseEnvsRoutineOpts struct {
	path string
	opts ParseOpts
	env  *v1alpha1.Environment
}

func parseEnvsRoutine(envsChan <-chan parseEnvsRoutineOpts) error {
	for req := range envsChan {
		_, env, err := ParseEnv(req.path, req.opts)
		if err != nil {
			return err
		}
		*req.env = *env
	}
	return nil
}
