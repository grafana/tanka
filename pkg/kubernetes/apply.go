package kubernetes

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/grafana/tanka/pkg/cli"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// ApplyOpts allow set additional parameters for the apply operation
type ApplyOpts client.ApplyOpts

// Apply receives a state object generated using `Reconcile()` and may apply it to the target system
func (k *Kubernetes) Apply(state manifest.List, opts ApplyOpts) error {
	alert := color.New(color.FgRed, color.Bold).SprintFunc()

	cluster := k.ctl.Info().Kubeconfig.Cluster
	context := k.ctl.Info().Kubeconfig.Context
	if !opts.AutoApprove {
		if err := cli.Confirm(
			fmt.Sprintf(`Applying to namespace '%s' of cluster '%s' at '%s' using context '%s'.`,
				alert(k.Spec.Namespace),
				alert(cluster.Name),
				alert(cluster.Cluster.Server),
				alert(context.Name),
			),
			"yes",
		); err != nil {
			return err
		}
	}
	return k.ctl.Apply(state, client.ApplyOpts(opts))
}
