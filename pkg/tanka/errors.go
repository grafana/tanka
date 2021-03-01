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
	return fmt.Sprintf("found multiple Environments in '%s': \n - %s", e.path, strings.Join(e.names, "\n - "))
}

// ErrParallel is an array of errors collected while processing in parallel
type ErrParallel struct {
	errors []error
}

func (e ErrParallel) Error() string {
	returnErr := fmt.Sprintf("Errors occured during parallel processing:\n\n")
	for _, err := range e.errors {
		returnErr = fmt.Sprintf("%s- %s\n\n", returnErr, err.Error())
	}
	return returnErr
}
