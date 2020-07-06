package tanka

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/term"
)

// Prune deletes all resources from the cluster, that are no longer present in
// Jsonnet. It uses the `tanka.dev/environment` label to identify those.
func Prune(baseDir string, mods ...Modifier) error {
	opts := parseModifiers(mods)

	// parse jsonnet, init k8s client
	p, err := load(baseDir, opts)
	if err != nil {
		return err
	}
	kube, err := p.connect()
	if err != nil {
		return err
	}
	defer kube.Close()

	// find orphaned resources
	orphaned, err := kube.Orphaned(p.Resources)
	if err != nil {
		return err
	}

	if len(orphaned) == 0 {
		fmt.Println("Nothing found to prune.")
		return nil
	}

	// print diff
	diff, err := kubernetes.StaticDiffer(false)(orphaned)
	if err != nil {
		// static diff can't fail normally, so unlike in apply, this is fatal
		// here
		return err
	}
	fmt.Print(term.Colordiff(*diff).String())

	// prompt for confirm
	if opts.apply.AutoApprove {
	} else if err := confirmPrompt("Pruning from", p.Env.Spec.Namespace, kube.Info()); err != nil {
		return err
	}

	// delete resources
	return kube.Delete(orphaned, opts.apply)
}
