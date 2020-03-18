package tanka

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
	"github.com/sh0rez/go-jsonnet/formatter"
)

// FormatOpts modify the behaviour of Format
type FormatOpts struct {
	// Excludes are a list of globs to exclude files while searching for Jsonnet
	// files
	Excludes []glob.Glob

	// OutFn receives the reformatted file and it's name. If left nil, the file
	// will be reformatted in place.
	OutFn OutFn

	// Whether to report changes in the return error. Make sure to set OutFn to
	// something non-nil if you don't want your files reformatted in-place.
	Test bool
}

// OutFn is a function that receives the reformatted file for further action,
// like persisting to disc
type OutFn func(name, content string) error

// Format takes files or directories, searches all Jsonnet files and
// reformats them.
func Format(fds []string, opts *FormatOpts) error {
	var paths []string
	for _, f := range fds {
		fs, err := findFiles(f, opts.Excludes)
		if err != nil {
			return errors.Wrap(err, "finding Jsonnet files")
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

	var changed []string
	for _, p := range paths {
		content, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		formatted, err := formatter.Format(p, string(content), formatter.DefaultOptions())
		if err != nil {
			return err
		}

		if string(content) != formatted {
			changed = append(changed, p)
		}

		if err := outFn(p, formatted); err != nil {
			return err
		}
	}

	if opts.Test && len(changed) != 0 {
		return ErrorNotFormatted{Files: changed}
	}

	return nil
}

// findFiles takes a file / directory and finds all files Jsonnet files
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
		Callback: func(path string, de *godirwalk.Dirent) error {
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

// ErrorNotFormatted means that one or more files need to be reformatted
type ErrorNotFormatted struct {
	// Files not properly formatted
	Files []string
}

func (e ErrorNotFormatted) Error() string {
	s := "The following files are not properly formatted:\n"
	for _, f := range e.Files {
		s += f + "\n"
	}
	return strings.TrimSuffix(s, "\n")
}
