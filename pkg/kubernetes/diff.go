package kubernetes

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/util"
)

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state manifest.List, opts DiffOpts) (*string, error) {
	if _, err := k.ctl.Resources(); err != nil {
		return nil, err
	}

	differ, err := k.differ(opts.Strategy)
	if err != nil {
		return nil, err
	}

	d, err := differ(state)
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

type separateOpts struct {
	namespaces map[string]bool
	resources  client.Resources
}

func separate(state manifest.List, defaultNs string, opts separateOpts) (live manifest.List, soon manifest.List) {
	soonNs := make(map[string]bool)
	for _, m := range state {
		if m.Kind() != "Namespace" {
			continue
		}
		soonNs[m.Metadata().Name()] = true
	}

	for _, m := range state {
		// non-namespaced always live
		if !opts.resources.Namespaced(m) {
			live = append(live, m)
			continue
		}

		// handle implicit default
		ns := m.Metadata().Namespace()
		if ns == "" {
			ns = defaultNs
		}

		// special case: namespace missing, BUT will be created during apply
		if !opts.namespaces[ns] && soonNs[ns] {
			soon = append(soon, m)
			continue
		}

		// everything else
		live = append(live, m)
	}

	return live, soon
}

// ErrorDiffStrategyUnknown occurs when a diff-strategy is requested that does
// not exist.
type ErrorDiffStrategyUnknown struct {
	Requested string
	differs   map[string]Differ
}

func (e ErrorDiffStrategyUnknown) Error() string {
	strats := []string{}
	for s := range e.differs {
		strats = append(strats, s)
	}
	return fmt.Sprintf("diff strategy `%s` does not exist. Pick one of: %v", e.Requested, strats)
}

func (k *Kubernetes) differ(override string) (Differ, error) {
	strategy := k.Spec.DiffStrategy
	if override != "" {
		strategy = override
	}

	d, ok := k.differs[strategy]
	if !ok {
		return nil, ErrorDiffStrategyUnknown{
			Requested: strategy,
			differs:   k.differs,
		}
	}

	return d, nil
}
