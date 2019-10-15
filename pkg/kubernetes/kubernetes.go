package kubernetes

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/objx"
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

// Manifest describes a single Kubernetes manifest
type Manifest map[string]interface{}

func (m Manifest) Kind() string {
	return m["kind"].(string)
}

func (m Manifest) Name() string {
	return m["metadata"].(map[string]interface{})["name"].(string)
}

func (m Manifest) Namespace() string {
	return m["metadata"].(map[string]interface{})["namespace"].(string)
}

// Reconcile receives the raw evaluated jsonnet as a marshaled json dict and
// shall return it reconciled as a state object of the target system
func (k *Kubernetes) Reconcile(raw map[string]interface{}, objectspecs ...*regexp.Regexp) (state []Manifest, err error) {
	docs, err := walkJSON(raw, "")
	out := make([]Manifest, 0, len(docs))
	if err != nil {
		return nil, errors.Wrap(err, "flattening manifests")
	}
	for _, d := range docs {
		m := objx.New(d)
		if k != nil && !m.Has("metadata.namespace") {
			m.Set("metadata.namespace", k.Spec.Namespace)
		}
		out = append(out, Manifest(m))
	}

	if len(objectspecs) > 0 {
		out = funk.Filter(out, func(i interface{}) bool {
			p := objectspec(i.(Manifest))
			for _, o := range objectspecs {
				if o.MatchString(strings.ToLower(p)) {
					return true
				}
			}
			return false
		}).([]Manifest)
	}

	sort.SliceStable(out, func(i int, j int) bool {
		if out[i].Kind() != out[j].Kind() {
			return out[i].Kind() < out[j].Kind()
		}
		return out[i].Name() < out[j].Name()
	})

	return out, nil
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

// Apply receives a state object generated using `Reconcile()` and may apply it to the target system
func (k *Kubernetes) Apply(state []Manifest, opts ApplyOpts) error {
	if k == nil {
		return ErrorMissingConfig{"apply"}
	}

	yaml, err := k.Fmt(state)
	if err != nil {
		return err
	}
	return k.client.Apply(yaml, k.Spec.Namespace, opts)
}

// DiffOpts allow to specify additional parameters for diff operations
type DiffOpts struct {
	Summarize bool
}

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state []Manifest, opts DiffOpts) (*string, error) {
	if k == nil {
		return nil, ErrorMissingConfig{"diff"}
	}
	yaml, err := k.Fmt(state)
	if err != nil {
		return nil, err
	}

	if k.Spec.DiffStrategy == "" {
		k.Spec.DiffStrategy = "native"
		if _, server, err := k.client.Version(); err == nil {
			if server.LessThan(semver.MustParse("1.13.0")) {
				k.Spec.DiffStrategy = "subset"
			}
		}
	}

	d, err := k.differs[k.Spec.DiffStrategy](yaml)
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

func objectspec(m Manifest) string {
	return fmt.Sprintf("%s/%s",
		m.Kind(),
		m.Name(),
	)
}
