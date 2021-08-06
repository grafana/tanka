package jsonnet

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/linter"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/jsonnet/native"
	"github.com/pkg/errors"
)

// LintOpts modifies the behaviour of Lint
type LintOpts struct {
	// Excludes are a list of globs to exclude files while searching for Jsonnet
	// files
	Excludes []glob.Glob

	// PrintNames causes all filenames to be printed
	PrintNames bool
}

// Lint takes a list of files and directories, processes them and prints
// out to stderr if there are linting warnings
func Lint(fds []string, opts *LintOpts) error {
	vm := jsonnet.MakeVM()
	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	var paths []string
	for _, f := range fds {
		fs, err := FindFiles(f, opts.Excludes)
		if err != nil {
			return errors.Wrap(err, "finding Jsonnet files")
		}
		paths = append(paths, fs...)
	}

	lintingFailed := false
	importedDirs := map[string]bool{}
	for _, file := range paths {
		if opts.PrintNames {
			log.Printf("Linting %s...\n", file)
		}

		dir := filepath.Dir(file)
		if _, ok := importedDirs[dir]; !ok {
			jpath, _, _, err := jpath.Resolve(dir)
			if err != nil {
				return errors.Wrap(err, "resolving JPATH")
			}

			vm.Importer(NewExtendedImporter(jpath))
			importedDirs[dir] = true
		}

		content, _ := ioutil.ReadFile(file)
		lintingFailed = lintingFailed || linter.LintSnippet(vm, os.Stderr, file, string(content))
	}

	if lintingFailed {
		return errors.New("Linting has failed for at least one file")
	}
	return nil
}
