package kubernetes

import (
	"fmt"

	"github.com/Masterminds/semver"

	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Kubernetes exposes methods to work with the Kubernetes orchestrator
type Kubernetes struct {
	Env v1alpha1.Environment

	// Client (kubectl)
	ctl client.Client

	// Diffing
	differs map[string]Differ // List of diff strategies
}

// Differ is responsible for comparing the given manifests to the cluster and
// returning differences (if any) in `diff(1)` format.
type Differ func(manifest.List) (*string, error)

// New creates a new Kubernetes with an initialized client
func New(env v1alpha1.Environment) (*Kubernetes, error) {
	// setup client
	ctl, err := client.New(env.Spec.APIServer)
	if err != nil {
		return nil, err
	}

	// setup diffing
	if env.Spec.DiffStrategy == "" {
		env.Spec.DiffStrategy = "native"

		if ctl.Info().ServerVersion.LessThan(semver.MustParse("1.13.0")) {
			env.Spec.DiffStrategy = "subset"
		}
	}

	k := Kubernetes{
		Env: env,
		ctl: ctl,
		differs: map[string]Differ{
			"native": ctl.DiffServerSide,
			"subset": SubsetDiffer(ctl),
		},
	}

	return &k, nil
}

// Close runs final cleanup
func (k *Kubernetes) Close() error {
	return k.ctl.Close()
}

// DiffOpts allow to specify additional parameters for diff operations
type DiffOpts struct {
	// Use `diffstat(1)` to create a histogram of the changes instead
	Summarize bool

	// Set the diff-strategy. If unset, the value set in the spec is used
	Strategy string
}

// Info about the client, etc.
func (k *Kubernetes) Info() client.Info {
	return k.ctl.Info()
}

func objectspec(m manifest.Manifest) string {
	return fmt.Sprintf("%s/%s",
		m.Kind(),
		m.Metadata().Name(),
	)
}
