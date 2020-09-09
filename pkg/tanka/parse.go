package tanka

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// loaded is the final result of all processing stages:
// 1. jpath.Resolve: Consruct import paths
// 2. parseSpec: load spec.json
// 3. evalJsonnet: evaluate Jsonnet to JSON
// 4. process.Process: post-processing
//
// Also connect() is provided to connect to the cluster for live operations
type loaded struct {
	Env       *v1alpha1.Config
	Resources manifest.List
}

// connect opens a connection to the backing Kubernetes cluster.
func (p *loaded) connect() (*kubernetes.Kubernetes, error) {
	env := *p.Env

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

// load runs all processing stages described at the Processed type
func load(dir string, opts Opts) (*loaded, error) {
	raw, env, err := eval(dir, opts.ExtCode)
	if err != nil {
		return nil, err
	}

	rec, err := process.Process(raw, *env, opts.Filters)
	if err != nil {
		return nil, err
	}

	return &loaded{
		Resources: rec,
		Env:       env,
	}, nil
}

// eval runs all processing stages describe at the Processed type apart from
// post-processing, thus returning the raw Jsonnet result.
func eval(dir string, extCode map[string]string) (raw interface{}, env *v1alpha1.Config, err error) {
	_, baseDir, rootDir, err := jpath.Resolve(dir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "resolving jpath")
	}

	env, err = parseSpec(baseDir, rootDir)
	if err != nil {
		return nil, nil, err
	}

	raw, err = evalJsonnet(baseDir, env, extCode)
	if err != nil {
		return nil, nil, errors.Wrap(err, "evaluating jsonnet")
	}

	return raw, env, nil
}

// parseEnv parses the `spec.json` of the environment and returns a
// *kubernetes.Kubernetes from it
func parseSpec(baseDir, rootDir string) (*v1alpha1.Config, error) {
	// name of the environment: relative path from rootDir
	name, _ := filepath.Rel(rootDir, baseDir)

	config, err := spec.ParseDir(baseDir, name)
	if err != nil {
		switch err.(type) {
		// the config includes deprecated fields
		case spec.ErrDeprecated:
			log.Println(err)
		// spec.json missing. we can still work with the default value
		case spec.ErrNoSpec:
			return config, nil
		// some other error
		default:
			return nil, errors.Wrap(err, "reading spec.json")
		}
	}

	return config, nil
}

// evalJsonnet evaluates the jsonnet environment at the given directory starting with
// `main.jsonnet`
func evalJsonnet(baseDir string, env *v1alpha1.Config, extCode map[string]string) (interface{}, error) {
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
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}
	return data, nil
}
