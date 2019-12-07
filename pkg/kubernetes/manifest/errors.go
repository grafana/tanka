package manifest

import (
	"fmt"
	"strings"
)

// SchemaError means that some expected fields were missing
type SchemaError struct {
	fields map[string]bool
	name   string
}

// Error returns the fields the manifest at the path is missing
func (s *SchemaError) Error() string {
	e := ""
	for k, missing := range s.fields {
		if !missing {
			continue
		}
		e += ", " + k
	}
	e = strings.TrimPrefix(e, ", ")
	return fmt.Sprintf("%smissing or invalid fields: %s", s.name, e)
}

func (s *SchemaError) add(field string) {
	if s.fields == nil {
		s.fields = make(map[string]bool)
	}
	s.fields[field] = true
}

// Missing returns whether a field is missing
func (s *SchemaError) Missing(field string) bool {
	return s.fields[field]
}

// WithName inserts a path into the error message
func (s *SchemaError) WithName(name string) *SchemaError {
	s.name = fmt.Sprintf("%s has ", name)
	return s
}
