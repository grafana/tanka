// Package tanka allows to use most of Tanka's features available on the
// command line programmatically as a Golang library. Keep in mind that the API
// is still experimental and may change without and signs of warnings while
// Tanka is still in alpha. Nevertheless, we try to avoid breaking changes.
package tanka

import (
	"io"
	"regexp"

	"github.com/grafana/tanka/pkg/kubernetes"
)

// parseModifiers parses all modifiers into an options struct
func parseModifiers(mods []Modifier) *options {
	o := &options{
		prune: kubernetes.PruneOpts{
			Prune:    true,
			AllKinds: false,
		},
	}

	for _, mod := range mods {
		mod(o)
	}

	o.prune.Force = o.apply.Force
	o.apply.PruneOpts = o.prune
	o.diff.PruneOpts = o.prune

	return o
}

type options struct {
	// io.Writer to write warnings and notices to
	wWarn io.Writer

	// target regular expressions to limit the working set
	targets []*regexp.Regexp

	// additional options for diff
	diff kubernetes.DiffOpts
	// additional options for apply
	apply kubernetes.ApplyOpts
	// additional options for pruning as part of diff/apply
	prune kubernetes.PruneOpts
}

// Modifier allow to influence the behavior of certain `tanka.*` actions. They
// are roughly equivalent to flags on the command line. See the `tanka.With*`
// functions for available options.
type Modifier func(*options)

// WithWarnWriter allows to provide a custom io.Writer that all warnings are
// written to
func WithWarnWriter(w io.Writer) Modifier {
	return func(opts *options) {
		opts.wWarn = w
	}
}

// WithTargets allows to submit regular expressions to limit the working set of
// objects (https://tanka.dev/targets/).
func WithTargets(t ...*regexp.Regexp) Modifier {
	return func(opts *options) {
		opts.targets = t
	}
}

// WithDiffStrategy allows to set the used diff strategy.
// An empty string is ignored.
func WithDiffStrategy(ds string) Modifier {
	return func(opts *options) {
		if ds != "" {
			opts.diff.Strategy = ds
		}
	}
}

// WithDiffSummarize enables summary mode, which invokes `diffstat(1)` on the
// set of changes to create an overview
func WithDiffSummarize(b bool) Modifier {
	return func(opts *options) {
		opts.diff.Summarize = b
	}
}

// WithApplyForce allows to invoke `kubectl apply` with the `--force` flag
func WithApplyForce(b bool) Modifier {
	return func(opts *options) {
		opts.apply.Force = b
	}
}

// WithApplyAutoApprove allows to skip the interactive approval
func WithApplyAutoApprove(b bool) Modifier {
	return func(opts *options) {
		opts.apply.AutoApprove = b
	}
}

// WithPrune enables pruning of orphaned resources
func WithPrune(b bool) Modifier {
	return func(opts *options) {
		opts.prune.Prune = b
	}
}

// WihtPruneAllKinds enables pruning all kinds, instead of only the most common
// ones. Please note that this is much slower.
func WithPruneAllKinds(b bool) Modifier {
	return func(opts *options) {
		opts.prune.AllKinds = b
	}
}
