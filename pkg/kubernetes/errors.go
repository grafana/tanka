package kubernetes

import (
	"errors"
)

var (
	// ErrorMissingConfig means that the `spec.json` is absent
	ErrorMissingConfig = errors.New("This operation requires additional configuration. Refer to https://tanka.dev/environments for instructions")
)
