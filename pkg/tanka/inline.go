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
	"github.com/rs/zerolog/log"
)

// InlineLoader loads an environment that is specified inline from within
// Jsonnet. The Jsonnet output is expected to hold a tanka.dev/Environment type,
// Kubernetes resources are expected at the `data` key of this very type
type InlineLoader struct{}

func (i *InlineLoader) Load(path string, opts LoaderOpts) (*v1alpha1.Environment, error) {
	return i.loadEnvironment(path, "", opts)
}

func (i *InlineLoader) Peek(path string, opts LoaderOpts) (*v1alpha1.Environment, error) {
	return i.loadEnvironment(path, "{data::{}}", opts)
}

func (i *InlineLoader) List(path string, opts LoaderOpts) ([]*v1alpha1.Environment, error) {
	opts.JsonnetOpts.EvalScript = fmt.Sprintf(MetadataEvalScript, opts.Name)
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
		env, err := inlineParse(path, raw)
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

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (i *InlineLoader) loadEnvironment(path, mixin string, opts LoaderOpts) (*v1alpha1.Environment, error) {
	// If the environment jsonnet eval path isn't already found, list the envs
	if opts.EvalExpression == "" {
		log.Debug().Str("name", opts.Name).Str("path", path).Str("mixin", mixin).Msg("No eval expression given when loading an environment, listing")
		envs, err := i.List(path, opts)
		if err != nil {
			return nil, err
		}
		if len(envs) > 1 {
			names := make([]string, 0, len(envs))
			for _, e := range envs {
				// If there's a full match on the given name, use this environment
				name := e.Metadata.Name
				if name == opts.Name {
					envs = []*v1alpha1.Environment{e}
					break
				}
				names = append(names, name)
			}
			if len(envs) > 1 {
				sort.Strings(names)
				return nil, ErrMultipleEnvs{path, opts.Name, names}
			}
		}

		if len(envs) == 0 {
			return nil, fmt.Errorf("found no matching environments; run 'tk env list %s' to view available options", path)
		}

		opts.EvalExpression = envs[0].Status.JsonnetExpression
	}

	opts.JsonnetOpts.EvalScript = PatternEvalScript(opts.EvalExpression + mixin)

	data, err := i.Eval(path, opts)
	if err != nil {
		return nil, fmt.Errorf("evaluating %s: %w", opts.JsonnetOpts.EvalScript, err)
	}

	env, err := inlineParse(path, data.(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	return env, nil
}

func inlineParse(path string, data map[string]interface{}) (*v1alpha1.Environment, error) {
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

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	env, err := spec.Parse(raw, namespace)
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
