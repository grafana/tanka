package kubernetes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/objx"

	"github.com/sh0rez/tanka/pkg/provider/util"
)

// Kubernetes provider bridges tanka to the Kubernetse orchestrator.
type Kubernetes struct {
	APIServer string `mapstructure:"apiserver"`
	Namespace string `mapstructure:"namespace"`
}

var client = Kubectl{}

type ErrorPrimitiveReached struct {
	path, key string
	primitive interface{}
}

func (e ErrorPrimitiveReached) Error() string {
	return fmt.Sprintf("recursion did not resolve in a valid Kubernetes object, "+
		"because one of `kind` or `apiVersion` is missing in path `.%s`."+
		" Found non-dict value `%s` of type `%T` instead.",
		e.path, e.key, e.primitive)
}

// Init makes the provider ready to be used
func (k *Kubernetes) Init() error {
	client.APIServer = k.APIServer
	return nil
}

// Reconcile receives the raw evaluated jsonnet as a marshaled json dict and
// shall return it reconciled as a state object of the target system
func (k *Kubernetes) Reconcile(raw map[string]interface{}) (state interface{}, err error) {
	docs, err := flattenManifest(raw, "")
	if err != nil {
		return nil, errors.Wrap(err, "flattening manifests")
	}
	for _, d := range docs {
		m := objx.New(d)
		m.Set("metadata.namespace", k.Namespace)
	}
	return docs, nil
}

// flattenManifest traverses deeply nested kubernetes manifest and extracts them into a flat map.
func flattenManifest(deep map[string]interface{}, path string) ([]map[string]interface{}, error) {
	flat := []map[string]interface{}{}

	for n, d := range deep {
		if n == "__ksonnet" {
			continue
		}
		if _, ok := d.(map[string]interface{}); !ok {
			return nil, ErrorPrimitiveReached{path, n, d}
		}
		m := objx.New(d)
		if m.Has("apiVersion") && m.Has("kind") {
			flat = append(flat, m)
		} else {
			f, err := flattenManifest(m, path+"."+n)
			if err != nil {
				return nil, err
			}
			flat = append(flat, f...)
		}
	}
	return flat, nil
}

// Fmt receives the state and reformats it to YAML Documents
func (k *Kubernetes) Fmt(state interface{}) (string, error) {
	return util.ShowYAMLDocs(state.([]map[string]interface{}))
}

// Apply receives a state object generated using `Reconcile()` and may apply it to the target system
func (k *Kubernetes) Apply(state interface{}) error {
	yaml, err := k.Fmt(state)
	if err != nil {
		return err
	}
	return client.Apply(yaml)
}

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state interface{}) (string, error) {
	yaml, err := k.Fmt(state)
	if err != nil {
		return "", err
	}
	return client.Diff(yaml)
}

// Cmd shall return a command to be available under `tk provider`
func (k *Kubernetes) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kubernetes",
		Short: "Kubernetes provider commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("not implemented")
		},
	}
}
