package tanka

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

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

func (i *InlineLoader) Load(path string, opts LoaderOpts) (*v1alpha1.Environment, error) {
	if opts.Name != "" {
		opts.JsonnetOpts.EvalScript = fmt.Sprintf(SingleEnvEvalScript, opts.Name)
	}

	data, err := i.Eval(path, opts)
	if err != nil {
		return nil, err
	}

	envs, err := extractEnvs(data)
	if err != nil {
		return nil, err
	}

	if len(envs) > 1 {
		names := make([]string, 0, len(envs))
		for _, e := range envs {
			// If there's a full match on the given name, use this environment
			if name := e.Metadata().Name(); name == opts.Name {
				envs = manifest.List{e}
				break
			} else {
				names = append(names, name)
			}
		}
		if len(envs) > 1 {
			sort.Strings(names)
			return nil, ErrMultipleEnvs{path, names}
		}
	}

	if len(envs) == 0 {
		return nil, fmt.Errorf("found no matching environments; run 'tk env list %s' to view available options", path)
	}

	// TODO: Re-serializing the entire env here. This is horribly inefficient
	envData, err := json.Marshal(envs[0])
	if err != nil {
		return nil, err
	}

	env, err := inlineParse(path, envData)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (i *InlineLoader) Peek(path string, opts LoaderOpts) (*v1alpha1.Environment, error) {
	opts.JsonnetOpts.EvalScript = MetadataEvalScript
	if opts.Name != "" {
		opts.JsonnetOpts.EvalScript = fmt.Sprintf(MetadataSingleEnvEvalScript, opts.Name)
	}
	return i.Load(path, opts)
}

func (i *InlineLoader) List(path string, opts LoaderOpts) ([]*v1alpha1.Environment, error) {
	opts.JsonnetOpts.EvalScript = MetadataEvalScript
	data, err := i.Eval(path, opts)
	if err != nil {
		return nil, err
	}

	list, err := extractEnvs(data)
	if err != nil {
		return nil, err
	}

	envs := make([]*v1alpha1.Environment, 0, len(list))
	for _, raw := range list {
		data, err := json.Marshal(raw)
		if err != nil {
			return nil, err
		}

		env, err := inlineParse(path, data)
		if err != nil {
			return nil, err
		}

		envs = append(envs, env)
	}

	return envs, nil
}

func (i *InlineLoader) Eval(path string, opts LoaderOpts) (interface{}, error) {
	// Can't provide env as extVar, as we need to evaluate Jsonnet first to know it
	opts.ExtCode.Set(environmentExtCode, `error "Using tk.env and std.extVar('tanka.dev/environment') is only supported for static environments. Directly access this data using standard Jsonnet instead."`)

	raw, err := evalJsonnet(path, opts.JsonnetOpts)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func inlineParse(path string, data []byte) (*v1alpha1.Environment, error) {
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

	env, err := spec.Parse(data, namespace)
	if err != nil {
		return nil, err
	}

	return env, nil
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
