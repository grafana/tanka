// Package tanka allows to use most of Tanka's features available on the
// command line programmatically as a Golang library. Keep in mind that the API
// is still experimental and may change without and signs of warnings while
// Tanka is still in alpha. Nevertheless, we try to avoid breaking changes.
package tanka

import (
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/process"
)

// parseModifiers parses all modifiers into an options struct
func parseModifiers(mods []Modifier) *options {
	o := &options{}
	for _, mod := range mods {
		mod(o)
	}
	return o
}

type options struct {
	// `std.extVar`
	extCode map[string]string

	// target regular expressions to limit the working set
	targets process.Matchers

	// additional options for diff
	diff kubernetes.DiffOpts

	// additional options for apply
	apply kubernetes.ApplyOpts
}

// Modifier allow to influence the behavior of certain `tanka.*` actions. They
// are roughly equivalent to flags on the command line. See the `tanka.With*`
// functions for available options.
type Modifier func(*options)

// WithExtCode allows to pass external variables (jsonnet code) to the VM
func WithExtCode(code map[string]string) Modifier {
	return func(opts *options) {
		opts.extCode = code
	}
}

// WithTargets allows to submit regular expressions to limit the working set of
// objects (https://tanka.dev/output-filtering/).
func WithTargets(t process.Matchers) Modifier {
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

// WithApplyValidate allows to invoke `kubectl apply` with the `--validate=false` flag
func WithApplyValidate(b bool) Modifier {
	return func(opts *options) {
		opts.apply.Validate = b
	}
}

// WithApplyAutoApprove allows to skip the interactive approval
func WithApplyAutoApprove(b bool) Modifier {
	return func(opts *options) {
		opts.apply.AutoApprove = b
	}
}
