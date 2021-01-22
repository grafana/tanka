package tanka

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/term"
)

// PruneOpts specify additional properties for the Prune action
type PruneOpts struct {
	Opts

	// AutoApprove skips the interactive approval
	AutoApprove bool
	// Force ignores any warnings kubectl might have
	Force bool
	// Also prune namespaces
	IncludeNamespaces bool
}

// Prune deletes all resources from the cluster, that are no longer present in
// Jsonnet. It uses the `tanka.dev/environment` label to identify those.
func Prune(baseDir string, opts PruneOpts) error {
	// parse jsonnet, init k8s client
	p, err := Load(baseDir, opts.Opts)
	if err != nil {
		return err
	}
	kube, err := p.Connect()
	if err != nil {
		return err
	}
	defer kube.Close()

	// find orphaned resources
	orphaned, err := kube.Orphaned(p.Resources, kubernetes.OrphanedOpts{opts.IncludeNamespaces})
	if err != nil {
		return err
	}

	if len(orphaned) == 0 {
		fmt.Println("Nothing found to prune.")
		return nil
	}

	// print diff
	diff, err := kubernetes.StaticDiffer(false, opts.IncludeNamespaces)(orphaned)
	if err != nil {
		// static diff can't fail normally, so unlike in apply, this is fatal
		// here
		return err
	}
	fmt.Print(term.Colordiff(*diff).String())

	// prompt for confirm
	if opts.AutoApprove {
	} else if err := confirmPrompt("Pruning from", p.Env.Spec.Namespace, kube.Info()); err != nil {
		return err
	}

	// delete resources
	return kube.Delete(orphaned, kubernetes.DeleteOpts{
		Force:             opts.Force,
		IncludeNamespaces: opts.IncludeNamespaces,
	})
}
