package jsonnet

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gobwas/glob"
	"github.com/google/go-jsonnet/linter"
	"github.com/grafana/tanka/pkg/jsonnet/implementations/goimpl"
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

	Out io.Writer
}

type lintWork struct {
	file   string
	jpaths []string
}

// Lint takes a list of files and directories, processes them and prints
// out to stderr if there are linting warnings
func Lint(fds []string, opts *LintOpts) error {
	if opts.Parallelism <= 0 {
		return errors.New("parallelism must be greater than 0")
	}

	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	// Resolve the jpath from each path the user passed (file or directory),
	// matching tk eval's entrypoint context, and reuse that jpath for every
	// file discovered under it. That preserves per-environment vendor
	// override priority and does not invent extra search paths that would
	// hide incorrect imports when a nested file is linted on its own.
	var work []lintWork
	for _, f := range fds {
		fs, err := FindFiles(f, opts.Excludes)
		if err != nil {
			return errors.Wrap(err, "finding Jsonnet files")
		}

		resolvePath, err := filepath.Abs(f)
		if err != nil {
			return errors.Wrapf(err, "getting absolute path for %s", f)
		}
		jpaths, _, _, err := jpath.Resolve(resolvePath, true)
		if err != nil {
			return errors.Wrapf(err, "getting jpath for %s", f)
		}

		for _, file := range fs {
			work = append(work, lintWork{file: file, jpaths: jpaths})
		}
	}

	type result struct {
		success bool
		output  string
	}
	workCh := make(chan lintWork, len(work))
	resultCh := make(chan result, len(work))
	lintWorker := func(workCh <-chan lintWork, resultCh chan result) {
		for item := range workCh {
			buf, success := lintWithRecover(item.file, item.jpaths)
			resultCh <- result{success: success, output: buf.String()}
		}
	}

	for i := 0; i < opts.Parallelism; i++ {
		go lintWorker(workCh, resultCh)
	}

	for _, item := range work {
		workCh <- item
	}
	close(workCh)

	lintingFailed := false
	for i := 0; i < len(work); i++ {
		result := <-resultCh
		lintingFailed = lintingFailed || !result.success
		if result.output != "" {
			fmt.Fprint(opts.Out, result.output)
		}
	}

	if lintingFailed {
		return errors.New("Linting has failed for at least one file")
	}
	return nil
}

func lintWithRecover(file string, jpaths []string) (buf bytes.Buffer, success bool) {
	file, err := filepath.Abs(file)
	if err != nil {
		fmt.Fprintf(&buf, "got an error getting the absolute path for %s: %v\n\n", file, err)
		return
	}

	log.Debug().Str("file", file).Msg("linting file")
	startTime := time.Now()
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(&buf, "caught a panic while linting %s: %v\n\n", file, err)
		}
		log.Debug().Str("file", file).Dur("duration_ms", time.Since(startTime)).Msg("linted file")
	}()

	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(&buf, "got an error reading file %s: %v\n\n", file, err)
		return
	}

	vm := goimpl.MakeRawVM(jpaths, nil, nil, 0)

	failed := linter.LintSnippet(vm, &buf, []linter.Snippet{{FileName: file, Code: string(content)}})
	return buf, !failed
}
