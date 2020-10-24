package tanka

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/google/go-jsonnet/formatter"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

// FormatOpts modify the behaviour of Format
type FormatOpts struct {
	// Excludes are a list of globs to exclude files while searching for Jsonnet
	// files
	Excludes []glob.Glob

	// OutFn receives the formatted file and it's name. If left nil, the file
	// will be formatted in place.
	OutFn OutFn

	// PrintNames causes all filenames to be printed
	PrintNames bool
}

// OutFn is a function that receives the formatted file for further action,
// like persisting to disc
type OutFn func(name, content string) error

// FormatFiles takes a list of files and directories, processes them and returns
// which files were formatted and perhaps an error.
func FormatFiles(fds []string, opts *FormatOpts) ([]string, error) {
	var paths []string
	for _, f := range fds {
		fs, err := findFiles(f, opts.Excludes)
		if err != nil {
			return nil, errors.Wrap(err, "finding Jsonnet files")
		}
		paths = append(paths, fs...)
	}

	// if nothing defined, default to save inplace
	outFn := opts.OutFn
	if outFn == nil {
		outFn = func(name, content string) error {
			return ioutil.WriteFile(name, []byte(content), 0644)
		}
	}

	// print each file?
	printFn := func(...interface{}) { return }
	if opts.PrintNames {
		printFn = func(i ...interface{}) { fmt.Println(i...) }
	}

	var changed []string
	for _, p := range paths {
		content, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}

		formatted, err := Format(p, string(content))
		if err != nil {
			return nil, err
		}

		if string(content) != formatted {
			printFn("fmt", p)
			changed = append(changed, p)
		} else {
			printFn("ok ", p)
		}

		if err := outFn(p, formatted); err != nil {
			return nil, err
		}
	}

	return changed, nil
}

// Format takes a file's name and contents and returns them in properly
// formatted. The file does not have to exist on disk.
func Format(filename string, content string) (string, error) {
	return formatter.Format(filename, content, formatter.DefaultOptions())
}

// findFiles takes a file / directory and finds all Jsonnet files
func findFiles(target string, excludes []glob.Glob) ([]string, error) {
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
