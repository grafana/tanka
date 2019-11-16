package spec

import "fmt"

type depreciation struct {
	old, new string
}

// ErrDeprecated is a non-fatal error that occurs when deprecated fields are
// used in the spec.json
type ErrDeprecated []depreciation

func (e ErrDeprecated) Error() string {
	buf := ""
	for _, d := range e {
		buf += fmt.Sprintf("Warning: `%s` is deprecated, use `%s` instead.\n", d.old, d.new)
	}
	return buf
}
