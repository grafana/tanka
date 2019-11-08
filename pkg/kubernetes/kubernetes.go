package kubernetes

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	funk "github.com/thoas/go-funk"

	"github.com/grafana/tanka/pkg/cli"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Kubernetes bridges tanka to the Kubernetse orchestrator.
type Kubernetes struct {
	ctl  client.Client
	Spec v1alpha1.Spec

	// Diffing
	differs map[string]Differ // List of diff strategies
}

type Differ func(client.Manifests) (*string, error)

// New creates a new Kubernetes with an initialized client
func New(s v1alpha1.Spec) (*Kubernetes, error) {
	ctl, err := client.New(s.APIServer)
	if err != nil {
		return nil, errors.Wrap(err, "creating client")
	}

	k := Kubernetes{
		Spec: s,
		ctl:  ctl,
	}

	k.differs = map[string]Differ{
		"native": k.ctl.DiffServerSide,
		"subset": SubsetDiffer(ctl),
	}
	return &k, nil
}

// Reconcile receives the raw evaluated jsonnet as a marshaled json dict and
// shall return it reconciled as a state object of the target system
func (k *Kubernetes) Reconcile(raw map[string]interface{}, objectspecs []*regexp.Regexp) (state client.Manifests, err error) {
	docs, err := walkJSON(raw, "")
	out := make(client.Manifests, 0, len(docs))
	if err != nil {
		return nil, errors.Wrap(err, "flattening manifests")
	}
	for _, d := range docs {
		m := objx.New(d)
		if k != nil && !m.Has("metadata.namespace") {
			m.Set("metadata.namespace", k.Spec.Namespace)
		}
		out = append(out, client.Manifest(m))
	}

	if len(objectspecs) > 0 {
		tmp := funk.Filter(out, func(i interface{}) bool {
			p := objectspec(i.(client.Manifest))
			for _, o := range objectspecs {
				if o.MatchString(strings.ToLower(p)) {
					return true
				}
			}
			return false
		}).([]client.Manifest)
		out = client.Manifests(tmp)
	}

	sort.SliceStable(out, func(i int, j int) bool {
		if out[i].Kind() != out[j].Kind() {
			return out[i].Kind() < out[j].Kind()
		}
		return out[i].Metadata().Name() < out[j].Metadata().Name()
	})

	return out, nil
}

type ApplyOpts client.ApplyOpts

// Apply receives a state object generated using `Reconcile()` and may apply it to the target system
func (k *Kubernetes) Apply(state client.Manifests, opts ApplyOpts) error {
	info, err := k.ctl.Info()
	if err != nil {
		return err
	}
	alert := color.New(color.FgRed, color.Bold).SprintFunc()

	if !opts.AutoApprove {
		if err := cli.Confirm(
			fmt.Sprintf(`Applying to namespace '%s' of cluster '%s' at '%s' using context '%s'.`,
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
	return k.ctl.Apply(state, client.ApplyOpts(opts))
}

// DiffOpts allow to specify additional parameters for diff operations
type DiffOpts struct {
	// Use `diffstat(1)` to create a histogram of the changes instead
	Summarize bool

	// Set the diff-strategy. If unset, the value set in the spec is used
	Strategy string
}

// Diff takes the desired state and returns the differences from the cluster
func (k *Kubernetes) Diff(state client.Manifests, opts DiffOpts) (*string, error) {
	strategy := k.Spec.DiffStrategy
	if opts.Strategy != "" {
		strategy = opts.Strategy
	}

	if strategy == "" {
		strategy = "native"

		info, err := k.ctl.Info()
		if err == nil && info.ServerVersion.LessThan(semver.MustParse("1.13.0")) {
			strategy = "subset"
		}
	}

	d, err := k.differs[strategy](state)
	switch {
	case err != nil:
		return nil, err
	case d == nil:
		return nil, nil
	}

	if opts.Summarize {
		return diffstat(*d)
	}

	return d, nil
}

func objectspec(m client.Manifest) string {
	return fmt.Sprintf("%s/%s",
		m.Kind(),
		m.Metadata().Name(),
	)
}
