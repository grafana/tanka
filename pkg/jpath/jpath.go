package jpath

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Resolve the given directory and resolves the jPath around it. This means it:
// - figures out the project root (the one with .jsonnetfile, vendor/ and lib/)
// - figures out the environments base directory (the one with the main.jsonnet)
//
// It then constructs a jPath with the base directory, vendor/ and lib/.
// This results in predictable imports, as it doesn't matter whether the user called
// called the command further down tree or not. A little bit like git.
func Resolve(workdir string, filename string) (path []string, base, root string) {
	root, err := findParentFile("jsonnetfile.json", workdir, "/")
	if err != nil {
		panic(err)
	}

	base, err = findParentFile(filename, workdir, root)
	if err != nil {
		panic(err)
	}

	return []string{
		base,
		filepath.Join(root, "vendor"),
		filepath.Join(root, "lib"),
	}, base, root
}

// findParentFile traverses the parent directory tree for the given `file`,
// starting from `start` and ending in `stop`. If the file is not found an error is returned.
func findParentFile(file, start, stop string) (string, error) {
	files, err := ioutil.ReadDir(start)
	if err != nil {
		return "", err
	}

	if dirContainsFile(files, file) {
		return start, nil
	} else if start == stop {
		return "", errors.New(file + " not found")
	}
	return findParentFile(file, filepath.Dir(start), stop)
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
