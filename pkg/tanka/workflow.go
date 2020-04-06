package tanka

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/term"
)

// Apply parses the environment at the given directory (a `baseDir`) and applies
// the evaluated jsonnet to the Kubernetes cluster defined in the environments
// `spec.json`.
func Apply(baseDir string, mods ...Modifier) error {
	opts := parseModifiers(mods)

	p, err := parse(baseDir, opts)
	if err != nil {
		return err
	}
	kube, err := p.newKube()
	if err != nil {
		return err
	}
	defer kube.Close()

	// show diff
	diff, err := kube.Diff(p.Resources, kubernetes.DiffOpts{})
	switch {
	case err != nil:
		// This is not fatal, the diff is not strictly required
		fmt.Println("Error diffing:", err)
	case diff == nil:
		tmp := "Warning: There are no differences. Your apply may not do anything at all."
		diff = &tmp
	}

	b := term.Colordiff(*diff)
	fmt.Print(b.String())

	// prompt for confirmation
	if err := applyPrompt(p.Env.Spec.Namespace, kube.Info()); err != nil {
		return err
	}

	return kube.Apply(p.Resources, opts.apply)
}

// applyPrompt asks the user for confirmation before apply
func applyPrompt(namespace string, info client.Info) error {
	alert := color.New(color.FgRed, color.Bold).SprintFunc()

	return term.Confirm(
		fmt.Sprintf(`Applying to namespace '%s' of cluster '%s' at '%s' using context '%s'.`,
			alert(namespace),
			alert(info.Kubeconfig.Cluster.Name),
			alert(info.Kubeconfig.Cluster.Cluster.Server),
			alert(info.Kubeconfig.Context.Name),
		),
		"yes",
	)
}

// Diff parses the environment at the given directory (a `baseDir`) and returns
// the differences from the live cluster state in `diff(1)` format. If the
// `WithDiffSummarize` modifier is used, a histogram created using `diffstat(1)`
// is returned instead.
// The cluster information is retrieved from the environments `spec.json`.
// NOTE: This function requires on `diff(1)`, `kubectl(1)` and perhaps `diffstat(1)`
func Diff(baseDir string, mods ...Modifier) (*string, error) {
	opts := parseModifiers(mods)

	p, err := parse(baseDir, opts)
	if err != nil {
		return nil, err
	}
	kube, err := p.newKube()
	if err != nil {
		return nil, err
	}
	defer kube.Close()

	return kube.Diff(p.Resources, opts.diff)
}

func Prune(baseDir string, mods ...Modifier) error {
	opts := parseModifiers(mods)

	// parse jsonnet, init k8s client
	p, err := parse(baseDir, opts)
	if err != nil {
		return err
	}
	kube, err := p.newKube()
	if err != nil {
		return err
	}
	defer kube.Close()

	// find orphaned resources
	// orphaned, err := kube.Orphaned(p.Resources)
	// if err != nil {
	// 	return err
	// }

	return kube.Prune(p.Resources)
}

// Show parses the environment at the given directory (a `baseDir`) and returns
// the list of Kubernetes objects.
// Tip: use the `String()` function on the returned list to get the familiar yaml stream
func Show(baseDir string, mods ...Modifier) (manifest.List, error) {
	opts := parseModifiers(mods)

	p, err := parse(baseDir, opts)
	if err != nil {
		return nil, err
	}

	return p.Resources, nil
}

// Eval returns the raw evaluated Jsonnet output (without any transformations)
func Eval(dir string, mods ...Modifier) (raw map[string]interface{}, err error) {
	opts := parseModifiers(mods)

	r, _, err := eval(dir, opts.extCode)
	if err != nil {
		return nil, err
	}
	return r, nil
}
