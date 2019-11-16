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

func (s *SchemaError) Error() string {
	e := ""
	for k, missing := range s.fields {
		if !missing {
			continue
		}
		e += ", " + k
	}
	e = strings.TrimPrefix(e, ", ")
	return fmt.Sprintf("%s missing or invalid fields: %s", s.name, e)
}

func (s *SchemaError) add(field string) {
	if s.fields == nil {
		s.fields = make(map[string]bool)
	}
	s.fields[field] = true
}

func (s *SchemaError) Missing(field string) bool {
	return s.fields[field]
}

func (s *SchemaError) WithName(name string) *SchemaError {
	s.name = name
	return s
}
