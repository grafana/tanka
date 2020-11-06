package tanka

import "fmt"

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
	path string
}

func (e ErrMultipleEnvs) Error() string {
	return fmt.Sprintf("found multiple Environments in '%s'", e.path)
}
