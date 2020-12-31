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
	Env       *v1alpha1.Environment
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
	_, env, err := ParseEnv(path, opts.JsonnetOpts)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, fmt.Errorf("no Tanka environment found")
	}

	rec, err := LoadManifests(env, opts.Filters)
	if err != nil {
		return nil, err
	}

	return &loaded{
		Resources: rec,
		Env:       env,
	}, nil
}

func LoadManifests(env *v1alpha1.Environment, filters process.Matchers) (manifest.List, error) {
	if env.Metadata.Name == "" {
		return nil, fmt.Errorf("Environment has no metadata.name set in %s", env.Metadata.Namespace)
	}

	if err := checkVersion(env.Spec.ExpectVersions.Tanka); err != nil {
		return nil, err
	}

	resources, err := process.Process(env.Data, *env, filters)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

// parseSpec parses the `spec.json` of the environment and returns a
// *kubernetes.Kubernetes from it
func parseSpec(path string) (*v1alpha1.Environment, error) {
	_, baseDir, rootDir, err := jpath.Resolve(path)
	if err != nil {
		return nil, err
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
			return nil, err
		}
	}

	return config, nil
}

// ParseEnv evaluates the jsonnet environment at the given file system path and
// optionally also returns and Environment object
func ParseEnv(path string, opts JsonnetOpts) (interface{}, *v1alpha1.Environment, error) {
	specEnv, err := parseSpec(path)
	if err != nil {
		switch err.(type) {
		case spec.ErrNoSpec:
			specEnv = nil
		default:
			return nil, nil, err
		}
	} else {
		// original behavior, if env has spec.json
		// then make env spec accessible through extCode
		jsonEnv, err := json.Marshal(specEnv)
		if err != nil {
			return nil, nil, errors.Wrap(err, "marshalling environment config")
		}
		// make a copy to prevent race condition in parallel execution
		injectCode := jsonnet.InjectedCode{}
		for k, v := range opts.ExtCode {
			injectCode.Set(k, v)
		}
		injectCode.Set(spec.APIGroup+"/environment", string(jsonEnv))
		opts.ExtCode = injectCode
	}

	raw, err := EvalJsonnet(path, opts)
	if err != nil {
		return nil, nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, nil, errors.Wrap(err, "unmarshalling data")
	}

	extractedEnvs, err := extractEnvironments(data)
	if _, ok := err.(process.ErrorPrimitiveReached); ok {
		if specEnv == nil {
			// if no environments or spec found, behave as jsonnet interpreter
			return data, nil, ErrNoEnv{path}
		}
	} else if err != nil {
		return nil, nil, err
	}

	var env *v1alpha1.Environment

	if len(extractedEnvs) > 1 {
		names := make([]string, 0)
		for _, exEnv := range extractedEnvs {
			names = append(names, exEnv.Metadata().Name())
		}
		return data, nil, ErrMultipleEnvs{path, names}
	} else if len(extractedEnvs) == 1 {
		marshalled, err := json.Marshal(extractedEnvs[0])
		if err != nil {
			return nil, nil, err
		}
		relpath, err := relativeEntrypoint(path)
		if err != nil {
			return nil, nil, err
		}
		env, err = spec.Parse(marshalled, relpath)
		if err != nil {
			return nil, nil, err
		}
		return data, env, nil
	} else if specEnv != nil {
		// if no environments found, fallback to original behavior
		specEnv.Data = data
		return data, specEnv, nil
	}
	// if no environments or spec found, behave as jsonnet interpreter
	return data, nil, ErrNoEnv{path}
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

// extractEnvironments filters out any Environment manifests
func extractEnvironments(data interface{}) (manifest.List, error) {
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

func relativeEntrypoint(path string) (string, error) {
	entrypoint, err := jpath.Entrypoint(path)
	if err != nil {
		return "", err
	}
	_, _, rootDir, err := jpath.Resolve(entrypoint)
	if err != nil {
		return "", errors.Wrap(err, "resolving jpath")
	}

	return filepath.Rel(rootDir, entrypoint)
}
