package kubernetes

import "fmt"

type ErrorNotFound struct {
	resource string
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf(`error from server (NotFound): secrets "%s" not found`, e.resource)
}

type ErrorMissingConfig struct {
	Verb string
}

func (e ErrorMissingConfig) Error() string {
	return fmt.Sprintf("%s requires additional configuration. Refer to https://tanka.dev/environments for that.", e.Verb)
}
