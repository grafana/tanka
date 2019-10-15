package tanka

import (
	"io"
	"regexp"

	"github.com/grafana/tanka/pkg/kubernetes"
)

func parseModifiers(mods []modifier) *options {
	o := &options{}
	for _, mod := range mods {
		mod(o)
	}
	return o
}

type options struct {
	wWarn io.Writer

	targets []*regexp.Regexp

	diff  kubernetes.DiffOpts
	apply kubernetes.ApplyOpts
}

type modifier func(*options)

// WithWarnWriter allows to provide a custom io.Writer that all warnings are
// written to
func WithWarnWriter(w io.Writer) modifier {
	return func(opts *options) {
		opts.wWarn = w
	}
}

// WithTargets allows to submit regular expressions to limit the working set of
// objects (https://tanka.dev/targets/).
func WithTargets(t ...*regexp.Regexp) modifier {
	return func(opts *options) {
		opts.targets = t
	}
}

// WithDiffStrategy allows to set the used diff strategy.
// An empty string is ignored.
func WithDiffStrategy(ds string) modifier {
	return func(opts *options) {
		if ds != "" {
			opts.diff.Strategy = ds
		}
	}
}

// WithDiffSummarize enables summary mode, which invokes `diffstat(1)` on the
// set of changes to create an overview
func WithDiffSummarize(b bool) modifier {
	return func(opts *options) {
		opts.diff.Summarize = b
	}
}

// WithApplyForce allows to invoke `kubectl apply` with the `--force` flag
func WithApplyForce(b bool) modifier {
	return func(opts *options) {
		opts.apply.Force = b
	}
}

// WithApplyAutoApprove allows to skip the interactive approval
func WithApplyAutoApprove(b bool) modifier {
	return func(opts *options) {
		opts.apply.AutoApprove = b
	}
}
