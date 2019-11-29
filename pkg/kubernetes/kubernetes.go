package kubernetes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/cli"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/util"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Kubernetes exposes methods to work with the Kubernetes orchestrator
type Kubernetes struct {
	Spec v1alpha1.Spec

	// Client (kubectl)
	ctl  client.Client
	info client.Info

	// Diffing
	differs map[string]Differ // List of diff strategies
}

// Differ is responsible for comparing the given manifests to the cluster and
// returning differences (if any) in `diff(1)` format.
type Differ func(manifest.List) (*string, error)

// New creates a new Kubernetes with an initialized client
func New(s v1alpha1.Spec) (*Kubernetes, error) {
	// setup client
	ctl, err := client.New(s.APIServer)
	if err != nil {
		return nil, errors.Wrap(err, "creating client")
	}

	// obtain information about the client (including versions)
	info, err := ctl.Info()
	if err != nil {
		return nil, err
	}

	// setup diffing
	if s.DiffStrategy == "" {
		s.DiffStrategy = "native"

		if info.ServerVersion.LessThan(semver.MustParse("1.13.0")) {
			s.DiffStrategy = "subset"
		}
	}

	k := Kubernetes{
		Spec: s,
		ctl:  ctl,
		info: *info,
		differs: map[string]Differ{
			"native": ctl.DiffServerSide,
			"subset": SubsetDiffer(ctl),
		},
	}

	return &k, nil
}

// ApplyOpts allow set additional parameters for the apply operation
type ApplyOpts client.ApplyOpts

// Apply receives a state object generated using `Reconcile()` and may apply it to the target system
func (k *Kubernetes) Apply(baseDir string, state manifest.List, opts ApplyOpts) error {
	info, err := k.ctl.Info()
	if err != nil {
		return err
	}
	alert := color.New(color.FgRed, color.Bold).SprintFunc()

	if !opts.AutoApprove {
		var message string
		if opts.Prune {
			message = `Applying to and pruning namespace '%s' of cluster '%s' at '%s' using context '%s'.`
		} else {
			message = `Applying to namespace '%s' of cluster '%s' at '%s' using context '%s'.`
		}
		if err := cli.Confirm(
			fmt.Sprintf(message,
				alert(k.Spec.Namespace),
				alert(info.Cluster.Get("name").MustStr()),
				alert(info.Cluster.Get("cluster.server").MustStr()),
				alert(info.Context.Get("name").MustStr()),
			),
			"yes",
		); err != nil {
			return err
		}
	}

	environmentLabel := TankaEnvironmentLabel + "=" + getEnvironmentLabel(baseDir)
	labels := []string{environmentLabel}
	return k.ctl.Apply(labels, state, client.ApplyOpts(opts))
}

// DiffOpts allow to specify additional parameters for diff operations
type DiffOpts struct {
	// Use `diffstat(1)` to create a histogram of the changes instead
	Summarize bool

	// Set the diff-strategy. If unset, the value set in the spec is used
	Strategy string
}

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(baseDir string, state manifest.List, opts DiffOpts) (*string, error) {
	strategy := k.Spec.DiffStrategy
	if opts.Strategy != "" {
		strategy = opts.Strategy
	}
	d, err := k.differs[strategy](state)
	if err != nil {
		return nil, err
	}

	environmentLabel := TankaEnvironmentLabel + "=" + getEnvironmentLabel(baseDir)
	labels := []string{environmentLabel}

	if opts.Summarize && d != nil {
		d, err = util.Diffstat(*d)
		if err != nil {
			return nil, err
		}
	}
	output, err := k.ctl.ApplyDryRun(labels, state)
	if err != nil {
		return nil, err
	}
	prunes := extractPrunes(output)
	if len(prunes) == 0 {
		return d, nil
	}
	lines := []string{}
	if d != nil {
		lines = append(lines, *d)
	}
	lines = append(lines, "--- Resources due for pruning:")
	lines = append(lines, prunes...)
	result := strings.Join(lines, "\n")
	return &result, nil
}

// Info about the client, etc.
func (k *Kubernetes) Info() client.Info {
	return k.info
}

func objectspec(m manifest.Manifest) string {
	return fmt.Sprintf("%s/%s",
		m.Kind(),
		m.Metadata().Name(),
	)
}

var prunedRE, _ = regexp.Compile("(.*) pruned \\(dry run\\)")

func extractPrunes(output string) []string {
	prunes := []string{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		match := prunedRE.FindStringSubmatch(line)
		if len(match) > 0 {
			prunes = append(prunes, "- "+match[1])
		}
	}
	return prunes
}
