package tanka

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/kubernetes"
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

	diff, err := kube.Diff(p.Resources, kubernetes.DiffOpts{
		PruneOpts: opts.prune,
	})
	if err != nil {
		return errors.Wrap(err, "diffing")
	}
	if diff == nil {
		tmp := "Warning: There are no differences. Your apply may not do anything at all."
		diff = &tmp
	}

	if opts.wWarn == nil {
		opts.wWarn = os.Stderr
	}
	fmt.Fprintln(opts.wWarn, *diff)

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

	return kube.Diff(p.Resources, opts.diff)
}

// Show parses the environment at the given directory (a `baseDir`) and returns
// the evaluated jsonnet in yaml form
func Show(baseDir string, mods ...Modifier) (string, error) {
	opts := parseModifiers(mods)

	p, err := parse(baseDir, opts)
	if err != nil {
		return "", err
	}

	return p.Resources.String(), nil
}
