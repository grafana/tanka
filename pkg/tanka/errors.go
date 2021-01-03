package tanka

import (
	"fmt"
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
	var listEnvs string
	for _, name := range e.names {
		listEnvs = listEnvs + fmt.Sprintf("  - %s \n", name)
	}
	return fmt.Sprintf("found multiple Environments in %s:\n%s", e.path, listEnvs)
}

// ErrParseParallel is an array of errors collected while parsing environments in parallel
type ErrParseParallel struct {
	Errors []error
}

func (e ErrParseParallel) Error() string {
	returnErr := fmt.Sprintf("Unable to parse selected Environments:\n\n")
	for _, err := range e.Errors {
		returnErr = fmt.Sprintf("%s- %s\n", returnErr, err.Error())
	}
	return returnErr
}
