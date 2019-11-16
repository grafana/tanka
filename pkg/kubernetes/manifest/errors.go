package manifest

import (
	"fmt"
	"strings"
)

// SchemaError means that some expected fields were missing
type SchemaError map[string]bool

func (s SchemaError) Error() string {
	e := ""
	for k, missing := range s {
		if !missing {
			continue
		}
		e += ", " + k
	}
	e = strings.TrimPrefix(e, ", ")
	return fmt.Sprint("expected fields missing: ", e)
}

func (s SchemaError) add(field string) {
	s[field] = true
}
