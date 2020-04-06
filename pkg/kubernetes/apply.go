package kubernetes

import (
	"fmt"
	"strings"
	"time"

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
	orphaned, err := k.Orphaned(state)
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

	live, err := k.ctl.GetByState(state)
	if err != nil {
		return nil, err
	}

	for _, m := range live {
		uids[m.Metadata().UID()] = true
	}

	return uids, nil
}

// Orphaned returns previously created resources that are missing from the
// local state. It uses UIDs to safely identify objects.
func (k *Kubernetes) Orphaned(state manifest.List) (manifest.List, error) {
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
		LabelEnvironment: k.Env.Metadata.NameLabel(),
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("done", time.Since(start))

	// filter unknown
	for _, m := range matched {
		// ignore known ones
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

		// record and skip from now on
		orphaned = append(orphaned, m)
		uids[m.Metadata().UID()] = true
	}

	return orphaned, nil
}
