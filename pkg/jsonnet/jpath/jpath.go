package jpath

import (
	"os"
	"path/filepath"
	"slices"
)

const DefaultEntrypoint = "main.jsonnet"

type WeightedJPath interface {
	Path() string
	Weight() int
}

type StaticallyWeightedJPath struct {
	weight int
	path   string
}

func NewStaticallyWeightedJPath(path string, weight int) *StaticallyWeightedJPath {
	return &StaticallyWeightedJPath{
		weight: weight,
		path:   path,
	}
}

func (jp *StaticallyWeightedJPath) Weight() int {
	return jp.weight
}

func (jp *StaticallyWeightedJPath) Path() string {
	return jp.path
}

// Resolve the given path and resolves the jPath around it. This means it:
// - figures out the project root (the one with .jsonnetfile, vendor/ and lib/)
// - figures out the environments base directory (usually the main.jsonnet)
//
// It then constructs a jPath with the base directory, vendor/ and lib/.
// This results in predictable imports, as it doesn't matter whether the user called
// called the command further down tree or not. A little bit like git.
func Resolve(path string, allowMissingBase bool, additionalJPaths []WeightedJPath) (jpath []string, base, root string, err error) {
	root, err = FindRoot(path)
	if err != nil {
		return nil, "", "", err
	}

	base, err = FindBase(path, root)
	if err != nil && allowMissingBase {
		base, err = FsDir(path)
		if err != nil {
			return nil, "", "", err
		}
	} else if err != nil {
		return nil, "", "", err
	}

	paths := make([]WeightedJPath, 0, 4+len(additionalJPaths))
	paths = append(paths, NewStaticallyWeightedJPath(filepath.Join(root, "vendor"), 300))
	paths = append(paths, NewStaticallyWeightedJPath(filepath.Join(base, "vendor"), 200))
	paths = append(paths, NewStaticallyWeightedJPath(filepath.Join(root, "lib"), 100))
	paths = append(paths, NewStaticallyWeightedJPath(base, 0))
	if additionalJPaths != nil {
		paths = append(paths, additionalJPaths...)
	}

	slices.SortStableFunc(paths, func(a, b WeightedJPath) int {
		return b.Weight() - a.Weight()
	})

	// TODO: Sort these paths with highest weight first
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		result = append(result, path.Path())
	}

	// The importer iterates through this list in reverse order
	return result, base, root, nil
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
		return DefaultEntrypoint, nil
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
