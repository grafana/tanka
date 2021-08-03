package process

import (
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"regexp"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

type FilterOptions struct {
	Exprs    Matchers
	Selector labels.Selector
}

type FilterOption func(options *FilterOptions)

// FilterWithOptions returns all elements of the list that match at least one expression
// and are not ignored
func FilterWithOptions(list manifest.List, options ...FilterOption) manifest.List {
	opts := &FilterOptions{}
	for _, o := range options {
		o(opts)
	}

	if len(opts.Exprs) == 0 && opts.Selector == nil {
		return list
	}

	out := make(manifest.List, 0, len(list))
	for _, m := range list {
		if len(opts.Exprs) > 0 && !opts.Exprs.MatchString(m.KindName()) {
			continue
		}
		if len(opts.Exprs) > 0 && opts.Exprs.IgnoreString(m.KindName()) {
			continue
		}

		if opts.Selector != nil && !opts.Selector.Matches(m.Metadata()) {
			continue
		}

		out = append(out, m)
	}
	return out
}

// Deprecated: Use FilterWithOptions instead
func Filter(list manifest.List, exprs Matchers) manifest.List {
	return FilterWithOptions(list, func(options *FilterOptions) {
		options.Exprs = exprs
	})
}

// Matcher is a single filter expression. The passed argument of Matcher is of the
// form `kind/name` (manifest.KindName())
type Matcher interface {
	MatchString(string) bool
}

// Ignorer is like matcher, but for explicitely ignoring resources
type Ignorer interface {
	IgnoreString(string) bool
}

// Matchers is a collection of multiple expressions.
// A matcher may also implement Ignorer to explicitely ignore fields
type Matchers []Matcher

// MatchString returns whether at least one expression (OR) matches the string
func (e Matchers) MatchString(s string) bool {
	b := false
	for _, exp := range e {
		b = b || exp.MatchString(s)
	}
	return b
}

func (e Matchers) IgnoreString(s string) bool {
	b := false
	for _, exp := range e {
		i, ok := exp.(Ignorer)
		if !ok {
			continue
		}
		b = b || i.IgnoreString(s)
	}
	return b
}

// RegExps is a helper to construct Matchers from regular expressions
func RegExps(rs []*regexp.Regexp) Matchers {
	xprs := make(Matchers, 0, len(rs))
	for _, r := range rs {
		xprs = append(xprs, r)
	}
	return xprs
}

func StrExps(strs ...string) (Matchers, error) {
	exps := make(Matchers, 0, len(strs))
	for _, raw := range strs {
		// trim exlamation mark, not supported by regex
		s := fmt.Sprintf(`(?i)^%s$`, strings.TrimPrefix(raw, "!"))

		// create regexp matcher
		var exp Matcher
		exp, err := regexp.Compile(s)
		if err != nil {
			return nil, ErrBadExpr{err}
		}

		// if negative (!), invert regex behaviour
		if strings.HasPrefix(raw, "!") {
			exp = NegMatcher{exp: exp}
		}
		exps = append(exps, exp)
	}
	return exps, nil
}

func MustStrExps(strs ...string) Matchers {
	exps, err := StrExps(strs...)
	if err != nil {
		panic(err)
	}
	return exps
}

// ErrBadExpr occurs when the regexp compiling fails
type ErrBadExpr struct {
	inner error
}

func (e ErrBadExpr) Error() string {
	return fmt.Sprintf("%s.\nSee https://tanka.dev/output-filtering/#regular-expressions for details on regular expressions.", strings.Title(e.inner.Error()))
}

// NexMatcher is a matcher that inverts the original behaviour
type NegMatcher struct {
	exp Matcher
}

func (n NegMatcher) MatchString(s string) bool {
	return true
}

func (n NegMatcher) IgnoreString(s string) bool {
	return n.exp.MatchString(s)
}
