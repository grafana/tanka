package jsonnet

import (
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/karrick/godirwalk"
)

// FindFiles takes a file / directory and finds all Jsonnet files
func FindFiles(target string, excludes []glob.Glob) ([]string, error) {
	// if it's a file, don't try to find children
	fi, err := os.Stat(target)
	if err != nil {
		return nil, err
	}
	if fi.Mode().IsRegular() {
		return []string{target}, nil
	}

	var files []string

	// godirwalk is faster than filepath.Walk, 'cause no os.Stat required
	err = godirwalk.Walk(target, &godirwalk.Options{
		Callback: func(rawPath string, de *godirwalk.Dirent) error {
			// Normalize slashes for Windows
			path := filepath.ToSlash(rawPath)

			if de.IsDir() {
				return nil
			}

			// excluded?
			for _, g := range excludes {
				if g.Match(path) {
					return nil
				}
			}

			// only .jsonnet or .libsonnet
			if ext := filepath.Ext(path); ext == ".jsonnet" || ext == ".libsonnet" {
				files = append(files, path)
			}
			return nil
		},
		// faster, no sort required
		Unsorted: true,
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
