package jpath

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var (
	// ErrorNoRoot means no rootDir was found in the parents
	ErrorNoRoot = errors.New("could not locate a tkrc.yaml or jsonnetfile.json in the parent directories, which is required to identify the project root.\nRefer to https://tanka.dev/directory-structure for more information")

	// ErrorNoBase means no baseDir was found in the parents
	ErrorNoBase = errors.New("could not locate a main.jsonnet in the parent directories, which is required as the entrypoint for the evaluation.\nRefer to https://tanka.dev/directory-structure for more information")
)

// ErrorFileNotFound means that the searched file was not found
type ErrorFileNotFound struct {
	filename string
}

func (e ErrorFileNotFound) Error() string {
	return e.filename + " not found"
}

// Resolve the given directory and resolves the jPath around it. This means it:
// - figures out the project root (the one with .jsonnetfile, vendor/ and lib/)
// - figures out the environments base directory (the one with the main.jsonnet)
//
// It then constructs a jPath with the base directory, vendor/ and lib/.
// This results in predictable imports, as it doesn't matter whether the user called
// called the command further down tree or not. A little bit like git.
func Resolve(workdir string) (path []string, base, root string, err error) {
	workdir, err = filepath.Abs(workdir)
	if err != nil {
		return nil, "", "", err
	}

	root, err = findRoot(workdir)
	if err != nil {
		return nil, "", "", err
	}

	base, err = FindParentFile("main.jsonnet", workdir, root)
	if err != nil {
		if _, ok := err.(ErrorFileNotFound); ok {
			return nil, "", "", ErrorNoBase
		}
		return nil, "", "", err
	}

	// The importer iterates through this list in reverse order
	return []string{
		filepath.Join(root, "vendor"),
		filepath.Join(base, "vendor"), // Look for a vendor folder in the base dir before using the root vendor
		filepath.Join(root, "lib"),
		base,
	}, base, root, nil
}

// findRoot searches for a rootDir by the following criteria:
// - tkrc.yaml is considered first, for a jb-independent way of marking the root
// - if it is not present (default), jsonnetfile.json is used.
func findRoot(start string) (dir string, err error) {
  // root path based on os
	var stop string
	if runtime.GOOS == "windows" {
	   stop = "C:\\"
	} else {
		stop = "/"
	}

	// try tkrc.yaml first
	root, err := FindParentFile("tkrc.yaml", start, stop)
	if err == nil {
		return root, nil
	}

	// otherwise use jsonnetfile.json
	root, err = FindParentFile("jsonnetfile.json", start, stop)
	if err != nil {
		if _, ok := err.(ErrorFileNotFound); ok {
			return "", ErrorNoRoot
		}
		return "", err
	}

	return root, nil
}

// FindParentFile traverses the parent directory tree for the given `file`,
// starting from `start` and ending in `stop`. If the file is not found an error is returned.
func FindParentFile(file, start, stop string) (string, error) {
	files, err := ioutil.ReadDir(start)
	if err != nil {
		return "", err
	}

	if dirContainsFile(files, file) {
		return start, nil
	} else if start == stop {
		return "", ErrorFileNotFound{file}
	}
	return FindParentFile(file, filepath.Dir(start), stop)
}

// dirContainsFile returns whether a file is included in a directory.
func dirContainsFile(files []os.FileInfo, filename string) bool {
	for _, f := range files {
		if f.Name() == filename {
			return true
		}
	}
	return false
}
