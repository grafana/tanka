package jpath

import "errors"

var (
	// ErrorNoRoot means no rootDir was found in the parents
	ErrorNoRoot = errors.New("could not locate a tkrc.yaml or jsonnetfile.json in the parent directories, which is required to identify the project root.\nRefer to https://tanka.dev/directory-structure for more information")

	// ErrorNoBase means no baseDir was found in the parents
	ErrorNoBase = errors.New("could not locate entrypoint (usually main.jsonnet) in the parent directories, which is required as the entrypoint for the evaluation.\nRefer to https://tanka.dev/directory-structure for more information")
)

// ErrorFileNotFound means that the searched file was not found
type ErrorFileNotFound struct {
	filename string
}

func (e ErrorFileNotFound) Error() string {
	return e.filename + " not found"
}
