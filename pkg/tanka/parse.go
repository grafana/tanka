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
	raw, env, err := eval(dir, opts.JsonnetOpts)
	if err != nil {
		return nil, err
	}

	if err := checkVersion(env.Spec.ExpectVersions.Tanka); err != nil {
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
func eval(dir string, opts jsonnet.Opts) (raw interface{}, env *v1alpha1.Config, err error) {
	_, baseDir, rootDir, err := jpath.Resolve(dir)
	if err != nil {
		return nil, nil, errors.Wrap(err, "resolving jpath")
	}

	env, err = parseSpec(baseDir, rootDir)
	if err != nil {
		return nil, nil, err
	}

	raw, err = evalJsonnet(baseDir, env, opts)
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

// evalJsonnet evaluates the jsonnet environment at the given directory
func evalJsonnet(baseDir string, env *v1alpha1.Config, opts jsonnet.Opts) (interface{}, error) {
	// make env spec accessible from Jsonnet
	jsonEnv, err := json.Marshal(env)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling environment config")
	}
	opts.ExtCode.Set(spec.APIGroup+"/environment", string(jsonEnv))

	// evaluate Jsonnet
	var raw string
	mainFile, err := jpath.GetEntrypoint(baseDir)
	if err != nil {
		return nil, err
	}

	if opts.EvalPattern != "" {
		evalScript := fmt.Sprintf("(import '%s').%s", mainFile, opts.EvalPattern)
		raw, err = jsonnet.Evaluate(mainFile, evalScript, opts)
		if err != nil {
			return nil, err
		}
	} else {
		raw, err = jsonnet.EvaluateFile(mainFile, opts)
		if err != nil {
			return nil, err
		}
	}
	// parse result
	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}
	return data, nil
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
