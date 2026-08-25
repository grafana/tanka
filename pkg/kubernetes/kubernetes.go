package kubernetes

import (
	"github.com/Masterminds/semver"

	"github.com/grafana/tanka/internal/telemetry"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

var tracer = telemetry.Tracer("kubernetes")

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

const strategyServer = "server"
const strategyNative = "native"
const strategyValidate = "validate"
const strategySubset = "subset"

// New creates a new Kubernetes with an initialized client
func New(env v1alpha1.Environment) (*Kubernetes, error) {
	// setup client
	var ctl *client.Kubectl
	var err error
	if len(env.Spec.ContextNames) < 1 {
		ctl, err = client.New(env.Spec.APIServer)
	} else {
		ctl, err = client.NewFromNames(env.Spec.ContextNames)
	}
	if err != nil {
		return nil, err
	}

	// setup diffing
	if env.Spec.DiffStrategy == "" {
		if env.Spec.ApplyStrategy == strategyServer {
			env.Spec.DiffStrategy = strategyServer
		} else {
			env.Spec.DiffStrategy = strategyNative
		}

		if ctl.Info().ServerVersion.LessThan(semver.MustParse("1.13.0")) {
			env.Spec.DiffStrategy = strategySubset
		}
	}

	k := Kubernetes{
		Env: env,
		ctl: ctl,
		differs: map[string]Differ{
			strategyNative:   ctl.DiffClientSide,
			strategyValidate: ctl.ValidateServerSide,
			strategyServer:   ctl.DiffServerSide,
			strategySubset:   SubsetDiffer(ctl),
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
	// Create a histogram of the changes instead
	Summarize bool
	// Find orphaned resources and include them in the diff
	WithPrune bool

	// Set the diff-strategy. If unset, the value set in the spec is used
	Strategy string
}

// Info about the client, etc.
func (k *Kubernetes) Info() client.Info {
	return k.ctl.Info()
}
