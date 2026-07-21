package process

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// regexMetaRe matches any regex metacharacter that would make a kind pattern
// non-literal (i.e. not a plain alphanumeric kind name like "StatefulSet").
var regexMetaRe = regexp.MustCompile(`[.+*?()\[\]{}|\\^$]`)

// Filter returns all elements of the list that match at least one expression
// and are not ignored
func Filter(list manifest.List, exprs Matchers) manifest.List {
	out := make(manifest.List, 0, len(list))
	for _, m := range list {
		if !exprs.MatchString(m.KindName()) {
			continue
		}
		if exprs.IgnoreString(m.KindName()) {
			continue
		}
		out = append(out, m)
	}
	return out
}

// Matcher is a single filter expression. The passed argument of Matcher is of the
// form `kind/name` (manifest.KindName())
type Matcher interface {
	MatchString(string) bool
}

// Ignorer is like matcher, but for explicitly ignoring resources
type Ignorer interface {
	IgnoreString(string) bool
}

// Matchers is a collection of multiple expressions.
// A matcher may also implement Ignorer to explicitly ignore fields
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

// KindsFor returns the set of resource kinds that positive matchers in the
// collection target. It is used to restrict cluster API queries to only the
// resource types relevant to the filter set, avoiding the cost of listing
// every resource type when only a few are needed.
//
// Returns (kinds, true) when all positive matchers specify a single literal
// kind (e.g. "statefulset/.*" → "StatefulSet"). Returns (nil, false) when
// kind restriction cannot be determined — either because no positive matchers
// exist, a matcher uses regex metacharacters in the kind position, or the
// matcher type is not a compiled regexp. Callers should query all resource
// kinds when ok is false.
func (e Matchers) KindsFor() ([]string, bool) {
	if len(e) == 0 {
		return nil, false
	}
	var kinds []string
	hasPositive := false
	for _, exp := range e {
		// Negative matchers (NegMatcher) implement Ignorer; skip them here
		// because they restrict results, not which kinds to query.
		if _, isIgnorer := exp.(Ignorer); isIgnorer {
			continue
		}
		hasPositive = true
		r, ok := exp.(*regexp.Regexp)
		if !ok {
			// Unknown matcher type — be conservative and query all kinds.
			return nil, false
		}
		kind, ok := kindFromFilterPattern(r.String())
		if !ok {
			return nil, false
		}
		kinds = append(kinds, kind)
	}
	if !hasPositive {
		// Only negative matchers present — no kind restriction.
		return nil, false
	}
	return kinds, true
}

// kindFromFilterPattern extracts the literal kind name from a regex pattern
// produced by StrExps (format: "(?i)^<user-input>$"). Returns ("", false)
// when the kind portion contains regex metacharacters or cannot be determined.
func kindFromFilterPattern(pattern string) (string, bool) {
	// Strip the anchors added by StrExps.
	s := strings.TrimPrefix(pattern, `(?i)^`)
	s = strings.TrimSuffix(s, `$`)

	kindPart, _, hasSlash := strings.Cut(s, "/")
	if !hasSlash {
		// No kind/name separator: pattern like "(?i)^statefulset$".
		// KindName() always contains a slash, so this would never match a
		// resource anyway. Return the pattern as-is so callers can include
		// the kind if it happens to be a literal, but it won't do harm.
		if regexMetaRe.MatchString(s) {
			return "", false
		}
		return s, true
	}

	if regexMetaRe.MatchString(kindPart) {
		return "", false
	}
	return kindPart, true
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
	caser := cases.Title(language.English)
	return fmt.Sprintf("%s.\nSee https://tanka.dev/output-filtering/#regular-expressions for details on regular expressions.", caser.String(e.inner.Error()))
}

// NexMatcher is a matcher that inverts the original behaviour
type NegMatcher struct {
	exp Matcher
}

func (n NegMatcher) MatchString(_ string) bool {
	return true
}

func (n NegMatcher) IgnoreString(s string) bool {
	return n.exp.MatchString(s)
}
