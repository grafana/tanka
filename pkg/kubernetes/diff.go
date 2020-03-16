package kubernetes

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/util"
	"github.com/pkg/errors"
)

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state manifest.List, opts DiffOpts) (*string, error) {
	// required for separating
	namespaces, err := k.ctl.Namespaces()
	if err != nil {
		return nil, errors.Wrap(err, "listing namespaces")
	}
	resources, err := k.ctl.Resources()
	if err != nil {
		return nil, errors.Wrap(err, "listing known api-resources")
	}

	// separate resources in groups
	//
	// soon: resources that have unmet dependencies that will be met during
	// apply. These will be diffed statically, because checking with the cluster
	// would cause an error
	//
	// live: all other resources
	live, soon := separate(state, k.Spec.Namespace, separateOpts{
		namespaces: namespaces,
		resources:  resources,
	})

	// differ for live resources
	liveDiff, err := k.differ(opts.Strategy)
	if err != nil {
		return nil, err
	}

	// reports all resources as new
	staticDiff := StaticDiffer(true)

	// run the diff
	d, err := multiDiff{
		{differ: liveDiff, state: live},
		{differ: staticDiff, state: soon},
	}.diff()

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

func StaticDiffer(create bool) Differ {
	return func(state manifest.List) (*string, error) {
		s := ""
		for _, m := range state {
			is, should := m.String(), ""
			if create {
				is, should = should, is
			}

			d, err := util.DiffStr(util.DiffName(m), is, should)
			if err != nil {
				return nil, err
			}
			s += d
		}

		if s == "" {
			return nil, nil
		}

		return &s, nil
	}
}

type multiDiff []struct {
	differ Differ
	state  manifest.List
}

func (m multiDiff) diff() (*string, error) {
	diff := ""
	for _, d := range m {
		s, err := d.differ(d.state)
		if err != nil {
			return nil, err
		}

		if s == nil {
			continue
		}
		diff += *s
	}

	if diff == "" {
		return nil, nil
	}
	return &diff, nil
}
