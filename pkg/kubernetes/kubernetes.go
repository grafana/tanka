package kubernetes

import (
	"bytes"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	yaml "gopkg.in/yaml.v2"

	"github.com/grafana/tanka/pkg/config/v1alpha1"
)

// Kubernetes bridges tanka to the Kubernetse orchestrator.
type Kubernetes struct {
	client Kubectl
	Spec   v1alpha1.Spec

	// Diffing
	differs map[string]Differ // List of diff strategies
}

type Differ func(yaml string) (string, error)

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

// Reconcile receives the raw evaluated jsonnet as a marshaled json dict and
// shall return it reconciled as a state object of the target system
func (k *Kubernetes) Reconcile(raw map[string]interface{}) (state []Manifest, err error) {
	docs, err := walkJSON(raw, "")
	out := make([]Manifest, 0, len(docs))
	if err != nil {
		return nil, errors.Wrap(err, "flattening manifests")
	}
	for _, d := range docs {
		m := objx.New(d)
		m.Set("metadata.namespace", k.Spec.Namespace)
		out = append(out, Manifest(m))
	}
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
func (k *Kubernetes) Apply(state []Manifest) error {
	yaml, err := k.Fmt(state)
	if err != nil {
		return err
	}
	return k.client.Apply(yaml, k.Spec.Namespace)
}

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state []Manifest) (string, error) {
	yaml, err := k.Fmt(state)
	if err != nil {
		return "", err
	}

	if k.Spec.DiffStrategy == "" {
		k.Spec.DiffStrategy = "native"
		if _, server, err := k.client.Version(); err == nil {
			if !server.GreaterThan(semver.MustParse("1.13.0")) {
				k.Spec.DiffStrategy = "subset"
			}
		}
	}

	return k.differs[k.Spec.DiffStrategy](yaml)
}
