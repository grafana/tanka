package kubernetes

import (
	"bytes"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	funk "github.com/thoas/go-funk"
	yaml "gopkg.in/yaml.v3"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Kubernetes bridges tanka to the Kubernetse orchestrator.
type Kubernetes struct {
	client Kubectl
	Spec   v1alpha1.Spec

	// Diffing
	differs map[string]Differ // List of diff strategies
}

type Differ func(yaml string) (*string, error)

// New creates a new Kubernetes
func New(s v1alpha1.Spec) *Kubernetes {
	k := Kubernetes{
		Spec: s,
	}
	k.client.APIServer = k.Spec.APIServer
	k.differs = map[string]Differ{
		"native": k.client.Diff,
		"subset": k.client.SubsetDiff,
	}
	return &k
}

// Compile extracts kubernetes Manifests from raw evaluated jsonnet <kind>/<name>,
// provided the manifests match the given regular expressions. It finds each manifest by
// recursively walking the jsonnet structure.
//
// In addition, we sort the manifests to ensure the order is consistent in each
// show/diff/apply cycle. This isn't necessary, but it does help users by producing
// consistent diffs.
func (k *Kubernetes) Compile(raw map[string]interface{}, kindNameMatchers []*regexp.Regexp) ([]Manifest, error) {
	manifests, err := walkJSON(raw)
	if err != nil {
		return nil, err
	}

	// If we have any matchers, we should filter all the manifests by matching against their
	// <kind>/<name> identifier.
	if len(kindNameMatchers) > 0 {
		manifests = funk.Filter(manifests, func(elem interface{}) bool {
			manifest := elem.(Manifest)
			kindName := strings.ToLower(manifest.KindName())
			for _, matcher := range kindNameMatchers {
				if matcher.MatchString(kindName) {
					return true
				}
			}

			return false
		}).([]Manifest)
	}

	// If we don't have a namespace, we want to set it to the default that is configured in
	// our kubernetes specification
	for idx, manifest := range manifests {
		if k != nil && manifest.Get("metadata.namespace").IsNil() {
			manifests[idx] = manifest.Set("metadata.namespace", k.Spec.Namespace)
		}
	}

	sort.SliceStable(manifests, func(i int, j int) bool {
		return manifests[i].GroupVersionKindName() < manifests[j].GroupVersionKindName()
	})

	return manifests, nil
}

// Fmt receives the state and reformats it to YAML Documents
func (k *Kubernetes) Fmt(state []Manifest) (string, error) {
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)

	for _, d := range state {
		if err := enc.Encode(d); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

// Apply receives a state object generated using `Compile()` and may apply it to the
// target system
func (k *Kubernetes) Apply(state []Manifest, opts ApplyOpts) error {
	if k == nil {
		return ErrorMissingConfig
	}

	yaml, err := k.Fmt(state)
	if err != nil {
		return err
	}
	return k.client.Apply(yaml, k.Spec.Namespace, opts)
}

// DiffOpts allow to specify additional parameters for diff operations
type DiffOpts struct {
	// Use `diffstat(1)` to create a histogram of the changes instead
	Summarize bool

	// Set the diff-strategy. If unset, the value set in the spec is used
	Strategy string
}

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state []Manifest, opts DiffOpts) (*string, error) {
	if k == nil {
		return nil, ErrorMissingConfig
	}
	yaml, err := k.Fmt(state)
	if err != nil {
		return nil, err
	}

	ds := k.Spec.DiffStrategy
	if opts.Strategy != "" {
		ds = opts.Strategy
	}
	if ds == "" {
		ds = "native"
		if _, server, err := k.client.Version(); err == nil {
			if server.LessThan(semver.MustParse("1.13.0")) {
				ds = "subset"
			}
		}
	}

	d, err := k.differs[ds](yaml)
	switch {
	case err != nil:
		return nil, err
	case d == nil:
		return nil, nil
	}

	if opts.Summarize {
		return diffstat(*d)
	}

	return d, nil
}
