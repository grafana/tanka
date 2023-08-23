package tanka

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const defaultParallelism = 8
const exportStrategyFanOut = "export-once-fan-out"

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

	if opts.Parallelism > len(envs) {
		log.Info().Int("parallelism", opts.Parallelism).Int("envs", len(envs)).Msg("Reducing parallelism to match number of environments")
		opts.Parallelism = len(envs)
	}

	for i := 0; i < opts.Parallelism; i++ {
		go parallelWorker(jobsCh, outCh)
	}

	fanoutEnvs := make(map[string][]string)
	for _, env := range envs {
		rootDir, err := jpath.FindRoot(env.Metadata.Namespace)
		if err != nil {
			return nil, errors.Wrap(err, "finding root")
		}
		path := filepath.Join(rootDir, env.Metadata.Namespace)

		if env.Spec.ExportStrategy == exportStrategyFanOut {
			fanoutEnvs[path] = append(fanoutEnvs[path], env.Metadata.Name)
			continue
		}

		o := opts.Opts
		// TODO: This is required because the map[string]string in here is not
		// concurrency-safe. Instead of putting this burden on the caller, find
		// a way to handle this inside the jsonnet package. A possible way would
		// be to make the jsonnet package less general, more tightly coupling it
		// to Tanka workflow thus being able to handle such cases
		o.JsonnetOpts = o.JsonnetOpts.Clone()
		o.Name = env.Metadata.Name
		jobsCh <- parallelJob{
			path: path,
			opts: o,
		}
	}
	for path, names := range fanoutEnvs {
		o := opts.Opts
		// TODO: This is required because the map[string]string in here is not
		// concurrency-safe. Instead of putting this burden on the caller, find
		// a way to handle this inside the jsonnet package. A possible way would
		// be to make the jsonnet package less general, more tightly coupling it
		// to Tanka workflow thus being able to handle such cases
		o.JsonnetOpts = o.JsonnetOpts.Clone()
		o.Name = exportStrategyFanOut
		jobsCh <- parallelJob{
			path:       path,
			fanoutEnvs: names,
			opts:       o,
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
	path       string
	fanoutEnvs []string
	opts       Opts
}

type parallelOut struct {
	env *v1alpha1.Environment
	err error
}

func parallelWorker(jobsCh <-chan parallelJob, outCh chan parallelOut) {
	for job := range jobsCh {
		log.Debug().Str("name", job.opts.Name).Str("path", job.path).Msg("Loading environment")
		startTime := time.Now()

		if job.opts.Name == exportStrategyFanOut {
			loadedEnvs, err := fanOut(job.path, job.opts.JsonnetOpts)
			if err != nil {
				err = fmt.Errorf("%s:\n %w", job.path, err)
			}
			if err == nil && len(job.fanoutEnvs) != len(loadedEnvs) {
				err = fmt.Errorf("%s:\n expected %d environments, got %d", job.path, len(job.fanoutEnvs), len(loadedEnvs))
			}
			// Always output the same number of environments as listed (otherwise, the job count will be off)
			if err != nil {
				for range job.fanoutEnvs {
					outCh <- parallelOut{env: nil, err: err}
				}
			} else {
				for _, env := range loadedEnvs {
					outCh <- parallelOut{env: env, err: err}
				}
			}

		} else {
			env, err := LoadEnvironment(job.path, job.opts)
			if err != nil {
				err = fmt.Errorf("%s:\n %w", job.path, err)
			}
			outCh <- parallelOut{env: env, err: err}
		}

		log.Debug().Str("name", job.opts.Name).Str("path", job.path).Dur("duration_ms", time.Since(startTime)).Msg("Finished loading environment")
	}
}

func fanOut(path string, opts jsonnet.Opts) ([]*v1alpha1.Environment, error) {
	raw, err := evalJsonnet(path, opts)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}

	list, err := extractEnvs(data)
	if err != nil {
		return nil, err
	}

	envs := make([]*v1alpha1.Environment, 0, len(list))
	for _, raw := range list {
		data, err := json.Marshal(raw)
		if err != nil {
			return nil, err
		}

		env, err := inlineParse(path, data)
		if err != nil {
			return nil, err
		}

		envs = append(envs, env)
	}

	return envs, nil

}
