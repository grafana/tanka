package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	funk "github.com/thoas/go-funk"
)

type ErrBadTargetExp struct {
	inner error
}

func (e ErrBadTargetExp) Error() string {
	return fmt.Sprintf("%s.\nSee https://tanka.dev/targets/#regular-expressions for details on regular expressions.", strings.Title(e.inner.Error()))
}

// CompileTargetExps compiles the regular expression for each target
func CompileTargetExps(strs []string) (exps []*regexp.Regexp, err error) {
	exps = make([]*regexp.Regexp, 0, len(strs))
	for _, raw := range strs {
		// surround the regular expression with start and end markers
		s := fmt.Sprintf(`(?i)^%s$`, raw)
		exp, err := regexp.Compile(s)
		if err != nil {
			return nil, ErrBadTargetExp{err}
		}
		exps = append(exps, exp)
	}
	return exps, nil
}

// MustCompileTargetExps is like CompileTargetExps but panics on error
// Meant for static code
func MustCompileTargetExps(strs ...string) (exps []*regexp.Regexp) {
	exps, err := CompileTargetExps(strs)
	if err != nil {
		panic(err)
	}
	return exps
}

// FilterTargets filters list to only include those matched by targets
func FilterTargets(list manifest.List, targets []*regexp.Regexp) manifest.List {
	if len(targets) == 0 {
		return list
	}

	tmp := funk.Filter(list, func(i interface{}) bool {
		p := objectspec(i.(manifest.Manifest))
		for _, t := range targets {
			if t.MatchString(p) {
				return true
			}
		}
		return false
	}).([]manifest.Manifest)

	return manifest.List(tmp)
}

func objectspec(m manifest.Manifest) string {
	return fmt.Sprintf("%s/%s",
		m.Kind(),
		m.Metadata().Name(),
	)
}
