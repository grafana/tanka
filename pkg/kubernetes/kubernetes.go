package kubernetes

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	yaml "gopkg.in/yaml.v2"
)

// Kubernetes bridges tanka to the Kubernetse orchestrator.
type Kubernetes struct {
	APIServer string `json:"apiServer"`
	Namespace string `json:"namespace"`
}

// Manifest describes a single Kubernetes manifest
type Manifest map[string]interface{}

var client = Kubectl{}

// Init prepares internals
func (k *Kubernetes) Init() error {
	client.APIServer = k.APIServer
	return nil
}

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
		m.Set("metadata.namespace", k.Namespace)
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
	return client.Apply(yaml)
}

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state []Manifest) (string, error) {
	yaml, err := k.Fmt(state)
	if err != nil {
		return "", err
	}
	return client.Diff(yaml)
}
