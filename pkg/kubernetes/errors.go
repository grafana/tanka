package kubernetes

import "fmt"

type ErrorNotFound struct {
	name string
	kind string
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf(`error from server (NotFound): %s "%s" not found`, e.kind, e.name)
}

type ErrorMissingConfig struct {
	verb string
}

func (e ErrorMissingConfig) Error() string {
	return fmt.Sprintf("%s requires additional configuration. Refer to https://tanka.dev/environments for that.", e.verb)
}
