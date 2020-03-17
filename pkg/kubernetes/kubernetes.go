package kubernetes

import (
	"fmt"

	"github.com/Masterminds/semver"

	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/util"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Kubernetes exposes methods to work with the Kubernetes orchestrator
type Kubernetes struct {
	Spec v1alpha1.Spec

	// Client (kubectl)
	ctl client.Client

	// Diffing
	differs map[string]Differ // List of diff strategies
}

// Differ is responsible for comparing the given manifests to the cluster and
// returning differences (if any) in `diff(1)` format.
type Differ func(manifest.List) (*string, error)

// New creates a new Kubernetes with an initialized client
func New(s v1alpha1.Spec) (*Kubernetes, error) {
	// setup client
	ctl, err := client.New(s.APIServer, s.Namespace)
	if err != nil {
		return nil, err
	}

	// setup diffing
	if s.DiffStrategy == "" {
		s.DiffStrategy = "native"

		if ctl.Info().ServerVersion.LessThan(semver.MustParse("1.13.0")) {
			s.DiffStrategy = "subset"
		}
	}

	k := Kubernetes{
		Spec: s,
		ctl:  ctl,
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

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state manifest.List, opts DiffOpts) (*string, error) {
	strategy := k.Spec.DiffStrategy
	if opts.Strategy != "" {
		strategy = opts.Strategy
	}

	d, err := k.differs[strategy](state)
	switch {
	case err != nil:
		return nil, err
	case d == nil:
		return nil, nil
	}

	if opts.Summarize {
		return util.Diffstat(*d)
	}

	return d, nil
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
