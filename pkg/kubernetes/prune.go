package kubernetes

import (
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/pkg/errors"
)

var defaultPruneKinds = []string{
	// core
	"ConfigMap",
	"Endpoints",
	"Namespace",
	"PersistentVolumeClaim",
	"PersistentVolume",
	"Pod",
	"ReplicationController",
	"Secret",
	"ServiceAccount",
	"Service",

	// jobs
	"DaemonSet",
	"Deployment",
	"ReplicaSet",
	"StatefulSet",

	// batch
	"Job",
	"CronJob",

	// networking
	"Ingress",

	// rbac
	"ClusterRole",
	"ClusterRoleBinding",
	"Role",
	"RoleBinding",
}

type PruneOpts struct {
	// Whether to remove orphaned objects from the cluster
	Prune bool

	// Check all kinds instead of only the most common ones
	AllKinds bool

	// Skip verification and force deleting
	Force bool
}

func (k *Kubernetes) prune(state manifest.List, opts PruneOpts) error {
	orphan, err := k.listOrphaned(state, opts.AllKinds)
	if err != nil {
		return errors.Wrap(err, "listing orphaned objects")
	}

	for _, o := range orphan {
		k.ctl.Delete(o.Metadata().Namespace(), o.Kind(), o.Metadata().Name(), client.DeleteOpts{
			Force: opts.Force,
		})
	}
	return nil
}

// listOrphaned returns all resources known to the cluster not present in
// Jsonnet
func (k *Kubernetes) listOrphaned(state manifest.List, all bool) (orphaned manifest.List, err error) {
	if k.orphaned != nil {
		return k.orphaned, nil
	}

	kinds := defaultPruneKinds
	if all {
		kinds, err = k.ctl.APIResources()
		if err != nil {
			return nil, errors.Wrap(err, "listing apiResources")
		}
	}

	results := make(chan (manifest.List))
	errs := make(chan (error))

	// list all objects matching our label
	for _, kind := range kinds {
		go k.parallelGetByLabels(kind, k.Env.Metadata.NameLabel(), results, errs)
	}

	var lastErr error
	for _ = range kinds {
		select {
		case list := <-results:
			for _, m := range list {
				// ComponentStatus resource is broken in Kubernetes versions
				// below 1.17, it will be returned even if the label does not
				// match. Ignoring it here is fine, as it is an internal object
				// type.
				if m.APIVersion() == "v1" && m.Kind() == "ComponentStatus" {
					continue
				}

				if state.Has(m) {
					continue
				}
				orphaned = append(orphaned, m)
			}
		case err := <-errs:
			lastErr = err
		}
	}
	close(results)
	close(errs)

	if lastErr != nil {
		return nil, lastErr
	}

	return orphaned, nil
}

func (k *Kubernetes) parallelGetByLabels(kind, envName string, r chan (manifest.List), e chan (error)) {
	list, err := k.ctl.GetByLabels("", kind, map[string]string{
		LabelEnvironment: envName,
	})
	if err != nil {
		e <- errors.Wrapf(err, "getting orphans of kind '%s':", kind)
	}
	r <- list
}
