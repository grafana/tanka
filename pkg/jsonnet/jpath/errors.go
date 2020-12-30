package jpath

import (
	"errors"
	"fmt"
)

var (
	// ErrorNoRoot means no rootDir was found in the parents
	ErrorNoRoot = errors.New(`Unable to identify the project root.
Tried to find 'tkrc.yaml' or 'jsonnetfile.json' in the parent directories.
Please refer to https://tanka.dev/directory-structure for more information`)

	// ErrorNoBase means no baseDir was found in the parents
	// ErrorNoBase = errors.New("could not locate entrypoint (usually main.jsonnet) in the parent directories, which is required as the entrypoint for the evaluation.\nRefer to https://tanka.dev/directory-structure for more information")
	// ErrorNoBase = errors.New(`Unable to identify the environments base directory.
// `)
)

type ErrorNoBase struct {
	filename string
}

func (e ErrorNoBase) Error() string {
	return fmt.Sprintf(`Unable to identify the environments base directory.
Tried to find '%s' in the parent directories.
Please refer to https://tanka.dev/directory-structure for more information`, e.filename)
}

// ErrorFileNotFound means that the searched file was not found
type ErrorFileNotFound struct {
	filename string
}

func (e ErrorFileNotFound) Error() string {
	return e.filename + " not found"
}
