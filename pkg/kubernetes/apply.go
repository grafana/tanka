package kubernetes

import (
	"fmt"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// ApplyOpts allow set additional parameters for the apply operation
type ApplyOpts client.ApplyOpts

// Apply receives a state object generated using `Reconcile()` and may apply it to the target system
func (k *Kubernetes) Apply(state manifest.List, opts ApplyOpts) error {
	return k.ctl.Apply(state, client.ApplyOpts(opts))
}

func (k *Kubernetes) Prune(state manifest.List) error {
	orphaned, err := k.listOrphaned(state)
	if err != nil {
		return err
	}

	for _, m := range orphaned {
		fmt.Println(m.Metadata().Namespace(), m.APIVersion(), m.Kind(), m.Metadata().Name())
	}
	return nil
}

func (k *Kubernetes) uids(state manifest.List) (map[string]bool, error) {
	uids := make(map[string]bool)

	for _, local := range state {
		ns := local.Metadata().Namespace()
		if ns == "" {
			ns = k.Env.Spec.Namespace
		}

		live, err := k.ctl.Get(ns, local.Kind(), local.Metadata().Name())
		if err != nil {
			return nil, err
		}
		uids[live.Metadata().UID()] = true
	}

	return uids, nil
}

// listOrphaned returns previously created resources that are missing from the
// local state. It uses UIDs to safely identify objects.
func (k *Kubernetes) listOrphaned(state manifest.List) (manifest.List, error) {
	apiResources, err := k.ctl.Resources()
	if err != nil {
		return nil, err
	}

	uids, err := k.uids(state)
	if err != nil {
		return nil, err
	}

	var orphaned manifest.List
	for _, r := range apiResources {
		if !strings.Contains(r.Verbs, "list") {
			continue
		}

		matched, err := k.ctl.GetByLabels("", r.FQN(), map[string]string{
			LabelEnvironment: k.Env.Metadata.NameLabel(),
		})
		if err != nil {
			return nil, err
		}

		// filter unknown using uids
		for _, m := range matched {
			if uids[m.Metadata().UID()] {
				continue
			}

			// ComponentStatus resource is broken in Kubernetes versions
			// below 1.17, it will be returned even if the label does not
			// match. Ignoring it here is fine, as it is an internal object
			// type.
			if m.APIVersion() == "v1" && m.Kind() == "ComponentStatus" {
				continue
			}

			orphaned = append(orphaned, m)

			// recorded. skip from now on
			uids[m.Metadata().UID()] = true
		}
	}

	return orphaned, nil
}
