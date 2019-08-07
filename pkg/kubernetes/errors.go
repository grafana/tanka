package kubernetes

import "fmt"

type ErrorNotFound struct {
	resource string
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf(`error from server (NotFound): secrets "%s" not found`, e.resource)
}
