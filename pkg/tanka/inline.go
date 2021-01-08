package tanka

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// InlineLoader loads an environment that is specified inline from within
// Jsonnet. The Jsonnet output is expected to hold a tanka.dev/Environment type,
// Kubernetes resources are expected at the `data` key of this very type
type InlineLoader struct{}

func (i *InlineLoader) Load(path string, opts JsonnetOpts) (*v1alpha1.Environment, error) {
	raw, err := EvalJsonnet(path, opts)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}

	envs, err := extractEnvs(data)
	if err != nil {
		return nil, err
	}

	if len(envs) > 1 {
		names := make([]string, 0, len(envs))
		for _, e := range envs {
			names = append(names, e.Metadata().Name())
		}
		return nil, ErrMultipleEnvs{path, names}
	}

	if len(envs) == 0 {
		return nil, fmt.Errorf("Found no environments in '%s'", path)
	}

	root, err := jpath.FindRoot(path)
	if err != nil {
		return nil, err
	}

	file, err := jpath.Entrypoint(path)
	if err != nil {
		return nil, err
	}

	namespace, err := filepath.Rel(root, file)
	if err != nil {
		return nil, err
	}

	// TODO: Re-serializing the entire env here. This is horribly inefficient
	envData, err := json.Marshal(envs[0])
	if err != nil {
		return nil, err
	}

	env, err := spec.Parse(envData, namespace)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (i *InlineLoader) Peek(path string, opts JsonnetOpts) (*v1alpha1.Environment, error) {
	opts.EvalScript = EnvsOnlyEvalScript
	return i.Load(path, opts)
}

// extractEnvs filters out any Environment manifests
func extractEnvs(data interface{}) (manifest.List, error) {
	// Scan for everything that looks like a Kubernetes object
	extracted, err := process.Extract(data)
	if err != nil {
		return nil, err
	}

	// Unwrap *List types
	if err := process.Unwrap(extracted); err != nil {
		return nil, err
	}

	out := make(manifest.List, 0, len(extracted))
	for _, m := range extracted {
		out = append(out, m)
	}

	// Extract only object of Kind: Environment
	return process.Filter(out, process.MustStrExps("Environment/.*")), nil
}
