package tanka

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

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
func parse(dir string, opts *options) (*ParseResult, error) {
	raw, env, err := eval(dir, opts.extCode)
	if err != nil {
		return nil, err
	}

	rec, err := kubernetes.Reconcile(raw, env.Spec, opts.targets)

	if err != nil {
		return nil, errors.Wrap(err, "reconciling")
	}

	if opts.applyLabels {
		applyLabels(rec, env)
	}

	return &ParseResult{
		Resources: rec,
		Env:       env,
	}, nil
}

func applyLabels(state manifest.List, env *v1alpha1.Config) {
	for _, manifest := range state {
		meta := manifest.Metadata()
		meta.SetLabel("app.kubernetes.io/managed-by", "tanka")
		meta.SetLabel("tanka.dev/environment", strings.Replace(env.Metadata.Name, "/", ".", -1))
	}
}

// eval returns the raw evaluated Jsonnet and the parsed env used for evaluation
func eval(dir string, extCode map[string]string) (raw map[string]interface{}, env *v1alpha1.Config, err error) {
	baseDir, env, err := loadDir(dir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "loading environment")
	}

	raw, err = evalJsonnet(baseDir, env, extCode)
	if err != nil {
		return nil, nil, errors.Wrap(err, "evaluating jsonnet")
	}

	return raw, env, nil
}

func loadDir(dir string) (baseDir string, env *v1alpha1.Config, err error) {
	_, baseDir, rootDir, err := jpath.Resolve(dir)
	if err != nil {
		return "", nil, errors.Wrap(err, "resolving jpath")
	}

	env, err = parseEnv(baseDir, rootDir)
	if err != nil {
		return "", nil, err
	}
	return baseDir, env, nil
}

// parseEnv parses the `spec.json` of the environment and returns a
// *kubernetes.Kubernetes from it
func parseEnv(baseDir, rootDir string) (*v1alpha1.Config, error) {
	// name of the environment: relative path from rootDir
	name, _ := filepath.Rel(rootDir, baseDir)

	config, err := spec.ParseDir(baseDir, name)
	if err != nil {
		switch err.(type) {
		// the config includes deprecated fields
		case spec.ErrDeprecated:
			log.Println(err)
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
	if err != nil {
		return nil, err
	}

	var dict map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &dict); err != nil {
		return nil, err
	}
	return dict, nil
}
