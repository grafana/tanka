package tanka

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// ParseResult contains the environments config and the manifests of this
// particular env
type ParseResult struct {
	Env       *v1alpha1.Config
	Resources manifest.List
}

func (p *ParseResult) newKube() (*kubernetes.Kubernetes, error) {
	kube, err := kubernetes.New(p.Env.Spec)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to Kubernetes")
	}
	return kube, nil
}

// parse loads the `spec.json`, evaluates the jsonnet and returns both, the
// kubernetes object and the reconciled manifests
func parse(baseDir string, opts *options) (*ParseResult, error) {
	_, baseDir, rootDir, err := jpath.Resolve(baseDir)
	if err != nil {
		return nil, errors.Wrap(err, "resolving jpath")
	}

	env, err := parseEnv(baseDir, rootDir, opts)
	if err != nil {
		return nil, err
	}

	raw, err := eval(baseDir)
	if err != nil {
		return nil, errors.Wrap(err, "evaluating jsonnet")
	}

	rec, err := kubernetes.Reconcile(raw, *env, opts.targets)

	if err != nil {
		return nil, errors.Wrap(err, "reconciling")
	}

	return &ParseResult{
		Resources: rec,
		Env:       env,
	}, nil
}

// Eval returns the raw evaluated Jsonnet output (without any transformations)
func Eval(baseDir string, mods ...Modifier) (raw map[string]interface{}, err error) {
	opts := parseModifiers(mods)

	r, _, err := eval(baseDir, opts)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// eval returns the raw evaluated Jsonnet and the parsed env used for evaluation
func eval(baseDir string, opts *options) (raw map[string]interface{}, env *v1alpha1.Config, err error) {
	if opts == nil {
		opts = &options{}
	}

	env, err = parseEnv(baseDir, opts)
	if err != nil {
		return nil, nil, err
	}

	raw, err = evalJsonnet(baseDir, env, opts.extCode)
	if err != nil {
		return nil, nil, errors.Wrap(err, "evaluating jsonnet")
	}

	return raw, env, nil
}

// parseEnv parses the `spec.json` of the environment and returns a
// *kubernetes.Kubernetes from it
func parseEnv(baseDir, rootDir string, opts *options) (*v1alpha1.Config, error) {
	// name of the environment: relative path from rootDir
	name, _ := filepath.Rel(rootDir, baseDir)

	config, err := spec.ParseDir(baseDir, name)
	if err != nil {
		switch err.(type) {
		// config is missing
		case viper.ConfigFileNotFoundError:
			return nil, kubernetes.ErrorMissingConfig
		// the config includes deprecated fields
		case spec.ErrDeprecated:
			if opts.wWarn == nil {
				opts.wWarn = os.Stderr
			}
			fmt.Fprint(opts.wWarn, err)
		// some other error
		default:
			return nil, errors.Wrap(err, "reading spec.json")
		}
	}

	return config, nil
}

// evalJsonnet evaluates the jsonnet environment at the given directory starting with
// `main.jsonnet`
func evalJsonnet(baseDir string, env *v1alpha1.Config, extCode map[string]string) (map[string]interface{}, error) {
	jsonEnv, err := json.Marshal(env)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling environment config")
	}

	ext := []jsonnet.Modifier{
		jsonnet.WithExtCode(spec.APIGroup+"/environment", string(jsonEnv)),
	}
	for k, v := range extCode {
		ext = append(ext, jsonnet.WithExtCode(k, v))
	}

	raw, err := jsonnet.EvaluateFile(
		filepath.Join(baseDir, "main.jsonnet"),
		ext...,
	)

	var dict map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &dict); err != nil {
		return nil, err
	}
	return dict, nil
}
