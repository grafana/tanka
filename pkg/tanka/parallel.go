package tanka

import (
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const defaultParallelism = 8

type parallelOpts struct {
	Opts
	Selector    labels.Selector
	Parallelism int
}

// parallelLoadEnvironments evaluates multiple environments in parallel
func parallelLoadEnvironments(envs []*v1alpha1.Environment, opts parallelOpts) ([]*v1alpha1.Environment, error) {
	jobsCh := make(chan parallelJob)
	outCh := make(chan parallelOut, len(envs))

	if opts.Parallelism <= 0 {
		opts.Parallelism = defaultParallelism
	}

	for i := 0; i < opts.Parallelism; i++ {
		go parallelWorker(jobsCh, outCh)
	}

	for _, env := range envs {
		o := opts.Opts

		// TODO: This is required because the map[string]string in here is not
		// concurrency-safe. Instead of putting this burden on the caller, find
		// a way to handle this inside the jsonnet package. A possible way would
		// be to make the jsonnet package less general, more tightly coupling it
		// to Tanka workflow thus being able to handle such cases
		o.JsonnetOpts = o.JsonnetOpts.Clone()

		o.Name = env.Metadata.Name
		path := env.Metadata.Namespace
		rootDir, err := jpath.FindRoot(path)
		if err != nil {
			return nil, errors.Wrap(err, "finding root")
		}
		jobsCh <- parallelJob{
			path: filepath.Join(rootDir, path),
			opts: o,
		}
	}
	close(jobsCh)

	var outenvs []*v1alpha1.Environment
	var errors []error
	for i := 0; i < len(envs); i++ {
		out := <-outCh
		if out.err != nil {
			errors = append(errors, out.err)
			continue
		}
		if opts.Selector == nil || opts.Selector.Empty() || opts.Selector.Matches(out.env.Metadata) {
			outenvs = append(outenvs, out.env)
		}
	}

	if len(errors) != 0 {
		return outenvs, ErrParallel{errors: errors}
	}

	return outenvs, nil
}

type parallelJob struct {
	path string
	opts Opts
}

type parallelOut struct {
	env *v1alpha1.Environment
	err error
}

func parallelWorker(jobsCh <-chan parallelJob, outCh chan parallelOut) {
	for job := range jobsCh {
		log.Debug().Str("name", job.opts.Name).Str("path", job.path).Msg("Loading environment")
		startTime := time.Now()

		env, err := LoadEnvironment(job.path, job.opts)
		if err != nil {
			err = fmt.Errorf("%s:\n %w", job.path, err)
		}
		outCh <- parallelOut{env: env, err: err}

		log.Debug().Str("name", job.opts.Name).Str("path", job.path).Dur("duration_ms", time.Since(startTime)).Msg("Finished loading environment")
	}
}
