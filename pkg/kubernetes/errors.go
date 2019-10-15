package kubernetes

import (
	"errors"
	"fmt"
)

type ErrorNotFound struct {
	name string
	kind string
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf(`error from server (NotFound): %s "%s" not found`, e.kind, e.name)
}

var (
	ErrorMissingConfig = errors.New("This operation requires additional configuration. Refer to https://tanka.dev/environments for instructions")
)
