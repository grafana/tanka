package manifest

import (
	"fmt"
	"strings"
)

// SchemaError means that some expected fields were missing
type SchemaError struct {
	Fields   map[string]bool
	Name     string
	Manifest Manifest
}

// Error returns the fields the manifest at the path is missing
func (s *SchemaError) Error() string {
	fields := make([]string, 0, len(s.Fields))
	for k, missing := range s.Fields {
		if !missing {
			continue
		}
		fields = append(fields, k)
	}

	if s.Name == "" {
		s.Name = "Resource"
	}

	msg := fmt.Sprintf("%s has missing or invalid fields: %s", s.Name, strings.Join(fields, ", "))

	if s.Manifest != nil {
		msg += fmt.Sprintf(":\n\n%s\n\nPlease check above object.", SampleString(s.Manifest.String()).Indent(2))
	}

	return msg
}

// SampleString is used for displaying code samples for error messages. It
// truncates the output to 10 lines
type SampleString string

func (s SampleString) String() string {
	lines := strings.Split(strings.TrimSpace(string(s)), "\n")
	truncate := len(lines) >= 10
	if truncate {
		lines = lines[0:10]
	}
	out := strings.Join(lines, "\n")
	if truncate {
		out += "\n..."
	}
	return out
}

func (s SampleString) Indent(n int) string {
	indent := strings.Repeat(" ", n)
	lines := strings.Split(s.String(), "\n")
	return indent + strings.Join(lines, "\n"+indent)
}
