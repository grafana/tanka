package kubernetes

import (
	"fmt"
	"strings"
	"time"

	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
)

// ApplyOpts allow set additional parameters for the apply operation
type ApplyOpts client.ApplyOpts

// Apply receives a state object generated using `Reconcile()` and may apply it to the target system
func (k *Kubernetes) Apply(state manifest.List, opts ApplyOpts) error {
	return k.ctl.Apply(state, client.ApplyOpts(opts))
}

// AnnoationLastApplied is the last-applied-configuration annotation used by kubectl
const AnnotationLastApplied = "kubectl.kubernetes.io/last-applied-configuration"

type OrphanedOpts struct {
	IncludeNamespaces bool
}

// Orphaned returns previously created resources that are missing from the
// local state. It uses UIDs to safely identify objects.
func (k *Kubernetes) Orphaned(state manifest.List, opts OrphanedOpts) (manifest.List, error) {
	if !k.Env.Spec.InjectLabels {
		return nil, fmt.Errorf(`spec.injectLabels is set to false in your spec.json. Tanka needs to add
a label to your resources to reliably detect which were removed from Jsonnet.
See https://tanka.dev/garbage-collection for more details.`)
	}

	apiResources, err := k.ctl.Resources()
	if err != nil {
		return nil, err
	}

	start := time.Now()
	fmt.Print("fetching UID's .. ")
	uids, err := k.uids(state)
	if err != nil {
		return nil, err
	}
	fmt.Println("done", time.Since(start))

	var orphaned manifest.List

	// join all kinds that support LIST into a comma separated string for
	// kubectl
	kinds := ""
	for _, r := range apiResources {
		if !strings.Contains(r.Verbs, "list") {
			continue
		}

		kinds += "," + r.FQN()
	}
	kinds = strings.TrimPrefix(kinds, ",")

	start = time.Now()
	fmt.Print("fetching previously created resources .. ")
	// get all resources matching our label
	matched, err := k.ctl.GetByLabels("", kinds, map[string]string{
		process.LabelEnvironment: k.Env.Metadata.NameLabel(),
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("done", time.Since(start))

	// filter unknown
	for _, m := range matched {
		if !opts.IncludeNamespaces && m.Kind() == "Namespace" {
			continue
		}

		// ignore known ones
		if uids[m.Metadata().UID()] {
			continue
		}

		// skip objects not created explicitely
		if _, ok := m.Metadata().Annotations()[AnnotationLastApplied]; !ok {
			continue
		}

		// record and skip from now on
		orphaned = append(orphaned, m)
		uids[m.Metadata().UID()] = true
	}

	return orphaned, nil
}

func (k *Kubernetes) uids(state manifest.List) (map[string]bool, error) {
	uids := make(map[string]bool)

	live, err := k.ctl.GetByState(state, client.GetByStateOpts{
		IgnoreNotFound: true,
	})
	if _, ok := err.(client.ErrorNothingReturned); ok {
		// return empty map of uids when kubectl returns nothing
		return uids, nil
	} else if err != nil {
		return nil, err
	}

	for _, m := range live {
		uids[m.Metadata().UID()] = true
	}

	return uids, nil
}
