package jsonnet

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gobwas/glob"
	"github.com/google/go-jsonnet/linter"
	"github.com/grafana/tanka/pkg/jsonnet/implementation/goimpl"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// LintOpts modifies the behaviour of Lint
type LintOpts struct {
	// Excludes are a list of globs to exclude files while searching for Jsonnet
	// files
	Excludes []glob.Glob

	// Parallelism determines the number of workers that will process files
	Parallelism int
}

// Lint takes a list of files and directories, processes them and prints
// out to stderr if there are linting warnings
func Lint(fds []string, opts *LintOpts) error {
	var paths []string
	for _, f := range fds {
		fs, err := FindFiles(f, opts.Excludes)
		if err != nil {
			return errors.Wrap(err, "finding Jsonnet files")
		}
		paths = append(paths, fs...)
	}

	type result struct {
		failed bool
		output string
	}
	fileCh := make(chan string, len(paths))
	resultCh := make(chan result, len(paths))
	lintWorker := func(fileCh <-chan string, resultCh chan result) {
		for file := range fileCh {
			buf := &bytes.Buffer{}
			var err error
			file, err = filepath.Abs(file)
			if err != nil {
				fmt.Fprintf(buf, "got an error getting the absolute path for %s: %v\n\n", file, err)
				resultCh <- result{failed: true, output: buf.String()}
				continue
			}

			log.Debug().Str("file", file).Msg("linting file")
			startTime := time.Now()

			jpaths, _, _, err := jpath.Resolve(file, true)
			if err != nil {
				fmt.Fprintf(buf, "got an error getting JPATH for %s: %v\n\n", file, err)
				resultCh <- result{failed: true, output: buf.String()}
				continue
			}
			vm := goimpl.MakeRawVM(jpaths, nil, nil, 0)

			content, _ := os.ReadFile(file)
			failed := linter.LintSnippet(vm, buf, []linter.Snippet{{FileName: file, Code: string(content)}})
			resultCh <- result{failed: failed, output: buf.String()}
			log.Debug().Str("file", file).Dur("duration_ms", time.Since(startTime)).Msg("linted file")
		}
	}

	for i := 0; i < opts.Parallelism; i++ {
		go lintWorker(fileCh, resultCh)
	}

	for _, file := range paths {
		fileCh <- file
	}
	close(fileCh)

	lintingFailed := false
	for i := 0; i < len(paths); i++ {
		result := <-resultCh
		lintingFailed = lintingFailed || result.failed
		if result.output != "" {
			fmt.Print(result.output)
		}
	}

	if lintingFailed {
		return errors.New("Linting has failed for at least one file")
	}
	return nil
}
