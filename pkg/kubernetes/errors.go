package kubernetes

import (
	"errors"
)

var (
	ErrorMissingConfig = errors.New("This operation requires additional configuration. Refer to https://tanka.dev/environments for instructions")
)
