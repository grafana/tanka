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
	kube, err := kubernetes.New(*p.Env)
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

// eval evaluates the jsonnet environment at the given directory starting with
// `main.jsonnet`
func eval(baseDir string) (map[string]interface{}, error) {
	raw, err := jsonnet.EvaluateFile(filepath.Join(baseDir, "main.jsonnet"))
	if err != nil {
		return nil, err
	}

	var dict map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &dict); err != nil {
		return nil, err
	}
	return dict, nil
}
