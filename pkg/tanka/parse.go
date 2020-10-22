package tanka

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// DEFAULT_DEV_VERSION is the placeholder version used when no actual semver is
// provided using ldflags
const DEFAULT_DEV_VERSION = "dev"

// CURRENT_VERSION is the current version of the running Tanka code
var CURRENT_VERSION = DEFAULT_DEV_VERSION

// loaded is the final result of all processing stages:
// TODO: remove or update this summary
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
func load(path string, opts Opts) (*loaded, error) {
	_, env, err := eval(path, opts.JsonnetOpts)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, fmt.Errorf("no Tanka environment found")
	}

	if err := checkVersion(env.Spec.ExpectVersions.Tanka); err != nil {
		return nil, err
	}

	rec, err := process.Process(env.Data, *env, opts.Filters)
	if err != nil {
		return nil, err
	}

	return &loaded{
		Resources: rec,
		Env:       env,
	}, nil
}

// parseSpec parses the `spec.json` of the environment and returns a
// *kubernetes.Kubernetes from it
func parseSpec(path string) (*v1alpha1.Config, error) {
	_, baseDir, rootDir, err := jpath.Resolve(path)
	if err != nil {
		return nil, errors.Wrap(err, "resolving jpath")
	}

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
			return config, err
		// some other error
		default:
			return nil, errors.Wrap(err, "reading spec.json")
		}
	}

	return config, nil
}

// eval evaluates the jsonnet environment at the given path
func eval(path string, opts jsonnet.Opts) (interface{}, *v1alpha1.Config, error) {
	var hasSpec bool
	specEnv, err := parseSpec(path)
	if err != nil {
		switch err.(type) {
		case spec.ErrNoSpec:
			hasSpec = false
		default:
			return nil, nil, errors.Wrap(err, "reading spec.json")
		}
	} else {
		hasSpec = true

		// original behavior, if env has spec.json
		// then make env spec accessible through extCode
		jsonEnv, err := json.Marshal(specEnv)
		if err != nil {
			return nil, nil, errors.Wrap(err, "marshalling environment config")
		}
		opts.ExtCode.Set(spec.APIGroup+"/environment", string(jsonEnv))
	}

	entrypoint, err := jpath.Entrypoint(path)
	if err != nil {
		return nil, nil, err
	}

	// evaluate Jsonnet
	var raw string
	if opts.EvalPattern != "" {
		evalScript := fmt.Sprintf("(import '%s').%s", entrypoint, opts.EvalPattern)
		raw, err = jsonnet.Evaluate(entrypoint, evalScript, opts)
		if err != nil {
			return nil, nil, errors.Wrap(err, "evaluating jsonnet")
		}
	} else {
		raw, err = jsonnet.EvaluateFile(entrypoint, opts)
		if err != nil {
			return nil, nil, errors.Wrap(err, "evaluating jsonnet")
		}
	}

	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, nil, errors.Wrap(err, "unmarshalling data")
	}

	if opts.EvalPattern != "" {
		// EvalPattern has no affinity with an environment, behave as jsonnet interpreter
		return data, nil, nil
	}

	var env *v1alpha1.Config
	switch data.(type) {
	case []interface{}:
		env = &v1alpha1.Config{}
		// data is array, do not try to unmarshal,
		// multiple envs currently unsupported
	default:
		if err := json.Unmarshal([]byte(raw), &env); err != nil {
			return nil, nil, errors.Wrap(err, "unmarshalling into v1alpha1.Config")
		}
	}

	// env is not a v1alpha1.Config, fallback to original behavior
	if env.Kind != "Environment" {
		if hasSpec {
			specEnv.Data = data
			// return env from spec.json
			return data, specEnv, nil
		} else {
			// No spec.json found, behave as jsonnet interpreter
			return data, nil, nil
		}
	}
	// return env AS IS
	return *env, env, nil
}

func checkVersion(constraint string) error {
	if constraint == "" {
		return nil
	}
	if CURRENT_VERSION == DEFAULT_DEV_VERSION {
		return nil
	}

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return fmt.Errorf("Parsing version constraint: '%w'. Please check 'spec.expectVersions.tanka'", err)
	}

	v, err := semver.NewVersion(CURRENT_VERSION)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid semantic version: '%w'.\nThis likely means your build of Tanka is broken, as this is a compile-time value. When in doubt, please raise an issue", CURRENT_VERSION, err)
	}

	if !c.Check(v) {
		return fmt.Errorf("Current version '%s' does not satisfy the version required by the environment: '%s'. You likely need to use another version of Tanka", CURRENT_VERSION, constraint)
	}

	return nil
}
