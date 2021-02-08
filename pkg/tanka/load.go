package tanka

import (
	"fmt"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/pkg/errors"
)

// environmentExtCode is the extCode ID `tk.env` uses underneath
// TODO: remove "import tk" and replace it with tanka-util
const environmentExtCode = spec.APIGroup + "/environment"

// Load loads the Environment at `path`. It automatically detects whether to
// load inline or statically
func Load(path string, opts Opts) (*LoadResult, error) {
	env, err := LoadEnvironment(path, opts)
	if err != nil {
		return nil, err
	}

	result, err := LoadManifests(env, opts.Filters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func LoadEnvironment(path string, opts Opts) (*v1alpha1.Environment, error) {
	loader, err := DetectLoader(path)
	if err != nil {
		return nil, err
	}

	env, err := loader.Load(path, LoaderOpts{opts.JsonnetOpts, opts.Name})
	if err != nil {
		return nil, err
	}

	return env, nil
}

func LoadManifests(env *v1alpha1.Environment, filters process.Matchers) (*LoadResult, error) {
	if err := checkVersion(env.Spec.ExpectVersions.Tanka); err != nil {
		return nil, err
	}

	processed, err := process.Process(*env, filters)
	if err != nil {
		return nil, err
	}

	return &LoadResult{Env: env, Resources: processed}, nil
}

// Peek loads the metadata of the environment at path. To get resources as well,
// use Load
func Peek(path string, opts Opts) (*v1alpha1.Environment, error) {
	loader, err := DetectLoader(path)
	if err != nil {
		return nil, err
	}

	return loader.Peek(path, LoaderOpts{opts.JsonnetOpts, opts.Name})
}

// List finds metadata of all environments at path that could possibly be
// loaded. List can be used to deal with multiple inline environments, by first
// listing them, choosing the right one and then only loading that one
func List(path string, opts Opts) ([]*v1alpha1.Environment, error) {
	loader, err := DetectLoader(path)
	if err != nil {
		return nil, err
	}

	return loader.List(path, LoaderOpts{opts.JsonnetOpts, opts.Name})
}

// Eval returns the raw evaluated Jsonnet
func Eval(path string, opts Opts) (interface{}, error) {
	loader, err := DetectLoader(path)
	if err != nil {
		return nil, err
	}

	return loader.Eval(path, LoaderOpts{opts.JsonnetOpts, opts.Name})
}

var registeredLoaders = []Loader{
	&StaticLoader{},
	&InlineLoader{},
}

// SetLoaders register custom loaders
func SetLoaders(loaders ...Loader) {
	registeredLoaders = loaders
}

// DetectLoader detects whether the environment is inline or static and picks
// the approriate loader
func DetectLoader(path string) (Loader, error) {
	_, base, err := jpath.Dirs(path)
	if err != nil {
		return nil, err
	}

	for i := range registeredLoaders {
		l := registeredLoaders[i]
		ok, err := l.Detect(base)
		if err != nil {
			return nil, err
		}
		if ok {
			return l, nil
		}
	}
	return nil, errors.New("no matched loader")
}

// Loader is an abstraction over the process of loading Environments
type Loader interface {
	// Detect check loader for detecting
	Detect(base string) (bool, error)

	// Load a single environment at path
	Load(path string, opts LoaderOpts) (*v1alpha1.Environment, error)

	// Peek only loads metadata and omits the actual resources
	Peek(path string, opts LoaderOpts) (*v1alpha1.Environment, error)

	// List returns metadata of all possible environments at path that can be
	// loaded
	List(path string, opts LoaderOpts) ([]*v1alpha1.Environment, error)

	// Eval returns the raw evaluated Jsonnet
	Eval(path string, opts LoaderOpts) (interface{}, error)
}

type LoaderOpts struct {
	JsonnetOpts
	Name string
}

type LoadResult struct {
	Env       *v1alpha1.Environment
	Resources manifest.List
}

func (l LoadResult) Connect() (*kubernetes.Kubernetes, error) {
	env := *l.Env

	// check env is complete
	s := ""
	if env.Spec.APIServer == "" {
		s += "  * spec.apiServer: No Kubernetes cluster endpoint specified"
	}
	if env.Spec.Namespace == "" {
		s += "  * spec.namespace: Default namespace missing"
	}
	if s != "" {
		return nil, fmt.Errorf("Your Environment's spec.json seems incomplete:\n%s\n\nPlease see https://tanka.dev/config for reference", s)
	}

	// connect client
	kube, err := kubernetes.New(env)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to Kubernetes")
	}

	return kube, nil
}
