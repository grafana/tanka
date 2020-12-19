package tanka

import (
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

const BASEDIR_INDICATOR = "main.jsonnet"

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
// ignores main.jsonnet if no environments found
func FindEnvironments(workdir string, selector labels.Selector) (envs []*v1alpha1.Environment, err error) {
	dirs, err := FindBaseDirs(workdir)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		_, env, err := ParseEnv(dir, ParseOpts{Evaluator: EnvsOnlyEvaluator})
		if err != nil {
			switch err.(type) {
			case ErrNoEnv:
				continue
			default:
				return nil, err
			}
		}

		if selector == nil || selector.Empty() || selector.Matches(env.Metadata) {
			envs = append(envs, env)
		}
	}

	return envs, nil
}
