package tanka

import (
	"fmt"
	"strings"
)

// ErrNoEnv means that the given jsonnet has no Environment object
// This must not be fatal, some operations work without
type ErrNoEnv struct {
	path string
}

func (e ErrNoEnv) Error() string {
	return fmt.Sprintf("unable to find an Environment in '%s'", e.path)
}

// ErrMultipleEnvs means that the given jsonnet has multiple Environment objects
type ErrMultipleEnvs struct {
	path  string
	names []string
}

func (e ErrMultipleEnvs) Error() string {
	return fmt.Sprintf("found multiple Environments (%s) in '%s'", strings.Join(e.names, ", "), e.path)
}

// ErrParseParallel is an array of errors collected while parsing environments in parallel
type ErrParseParallel struct {
	errors []error
}

func (e ErrParseParallel) Error() string {
	returnErr := fmt.Sprintf("Unable to parse selected Environments:\n\n")
	for _, err := range e.errors {
		returnErr = fmt.Sprintf("%s- %s\n", returnErr, err.Error())
	}
	return returnErr
}
