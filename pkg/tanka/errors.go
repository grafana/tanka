package tanka

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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

// ErrParseEnvs is an array of errors collected while parsing environments in parallel
type ErrParseEnvs struct {
	errors []error
}

func (e ErrParseEnvs) Error() string {
	returnErr := errors.New("Unable to parse selected Environments")
	for _, err := range e.errors {
		returnErr = errors.Wrap(returnErr, err.Error())
	}
	return returnErr.Error()
}
