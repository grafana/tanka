package tanka

import (
	"fmt"

	"github.com/grafana/tanka/pkg/cli"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Apply parses the environment at the given directory (a `baseDir`) and applies
// the evaluated jsonnet to the Kubernetes cluster defined in the environments
// `spec.json`.
// NOTE: This function prints on screen in default configuration.
// Use the `WithWarnWriter` modifier to change that. The `WithApply*` modifiers
// may be used to further influence the behavior.
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

	diff, err := kube.Diff(p.Resources, kubernetes.DiffOpts{})
	switch {
	case err != nil:
		// This is not fatal, the diff is not strictly required
		fmt.Println("Error diffing:", err)
	case diff == nil:
		tmp := "Warning: There are no differences. Your apply may not do anything at all."
		diff = &tmp
	}

	b := cli.Colordiff(*diff)
	fmt.Print(b.String())

	return kube.Apply(p.Resources, opts.apply)
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
