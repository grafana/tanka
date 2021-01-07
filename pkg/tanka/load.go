package tanka

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/pkg/errors"
)

// Load loads the Environment at `path`. It automatically detects whether to
// load inline or statically
func Load(path string, opts Opts) (*LoadResult, error) {
	_, base, err := jpath.Dirs(path)
	if err != nil {
		return nil, err
	}

	loader, err := detectLoader(base)
	if err != nil {
		return nil, err
	}

	env, err := loader.Load(path, opts.JsonnetOpts)
	if err != nil {
		return nil, err
	}

	if err := checkVersion(env.Spec.ExpectVersions.Tanka); err != nil {
		return nil, err
	}

	processed, err := process.Process(*env, opts.Filters)
	if err != nil {
		return nil, err
	}

	return &LoadResult{Env: env, Resources: processed}, nil
}

// Peek loads the metadata of the environment at path. To get resources as well,
// use Load
func Peek(path string, opts JsonnetOpts) (*v1alpha1.Environment, error) {
	_, base, err := jpath.Dirs(path)
	if err != nil {
		return nil, err
	}

	loader, err := detectLoader(base)
	if err != nil {
		return nil, err
	}

	return loader.Peek(path, opts)
}

// detectLoader detects whether the environment is inline or static and picks
// the approriate loader
func detectLoader(base string) (Loader, error) {
	// check if spec.json exists
	_, err := os.Stat(filepath.Join(base, spec.Specfile))
	if os.IsNotExist(err) {
		return &InlineLoader{}, nil
	} else if err != nil {
		return nil, err
	}

	return &StaticLoader{}, nil
}

// Loader is an abstraction over the process of loading Environments
type Loader interface {
	// Load the environment with path
	Load(path string, opts JsonnetOpts) (*v1alpha1.Environment, error)

	// Peek only loads metadata and omits the actual resources
	Peek(path string, opts JsonnetOpts) (*v1alpha1.Environment, error)
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
