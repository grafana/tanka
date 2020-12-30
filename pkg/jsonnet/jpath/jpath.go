package jpath

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const DEFAULT_ENTRYPOINT = "main.jsonnet"

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

// Resolve the given path and resolves the jPath around it. This means it:
// - figures out the project root (the one with .jsonnetfile, vendor/ and lib/)
// - figures out the environments base directory (usually the main.jsonnet)
//
// It then constructs a jPath with the base directory, vendor/ and lib/.
// This results in predictable imports, as it doesn't matter whether the user called
// called the command further down tree or not. A little bit like git.
func Resolve(path string) (jpath []string, base, root string, err error) {
	root, err = FindRoot(path)
	if err != nil {
		return nil, "", "", err
	}

	base, err = FindBase(path, root)
	if err != nil {
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

// FindRoot returns the absolute path of the project root, being the directory
// that directly holds `tkrc.yaml` if it exists, otherwise the directory that
// directly holds `jsonnetfile.json`
func FindRoot(path string) (dir string, err error) {
	start, err := FsDir(path)
	if err != nil {
		return "", err
	}

	// root path based on os
	stop := "/"
	if runtime.GOOS == "windows" {
		stop = filepath.VolumeName(start) + "\\"
	}

	// try tkrc.yaml first
	root, err := FindParentFile("tkrc.yaml", start, stop)
	if err == nil {
		return root, nil
	}

	// otherwise use jsonnetfile.json
	root, err = FindParentFile("jsonnetfile.json", start, stop)
	if _, ok := err.(ErrorFileNotFound); ok {
		return "", ErrorNoRoot
	} else if err != nil {
		return "", err
	}

	return root, nil
}

// FindBase returns the absolute path of the environments base directory, the
// one which directly holds the entrypoint file.
func FindBase(path string, root string) (string, error) {
	dir, err := FsDir(path)
	if err != nil {
		return "", err
	}

	filename, err := Filename(path)
	if err != nil {
		return "", err
	}

	base, err := FindParentFile(filename, dir, root)

	if _, ok := err.(ErrorFileNotFound); ok {
		return "", ErrorNoBase
	} else if err != nil {
		return "", err
	}

	return base, nil
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

// FsDir returns the most inner directory of path, as reported by the local
// filesystem
func FsDir(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		return path, nil
	}

	return filepath.Dir(path), nil
}

// Filename returns the name of the entrypoint file.
// It DOES NOT return an absolute path, only a plain name like "main.jsonnet"
// To obtain an absolute path, use Entrypoint() instead.
func Filename(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		return DEFAULT_ENTRYPOINT, nil
	}

	return filepath.Base(fi.Name()), nil

}

// Entrypoint returns the absolute path of the environments entrypoint file (the
// one passed to jsonnet.EvaluateFile)
func Entrypoint(path string) (string, error) {
	root, err := FindRoot(path)
	if err != nil {
		return "", err
	}

	base, err := FindBase(path, root)
	if err != nil {
		return "", err
	}

	filename, err := Filename(path)
	if err != nil {
		return "", err
	}

	return filepath.Join(base, filename), nil
}
