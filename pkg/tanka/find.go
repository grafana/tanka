package tanka

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/labels"
)

// FindOpts are optional arguments for FindEnvs
type FindOpts struct {
	JsonnetOpts
	Selector    labels.Selector
	Parallelism int
}

// FindEnvs returns metadata of all environments recursively found in 'path'.
// Each directory is tested and included if it is a valid environment, either
// static or inline. If a directory is a valid environment, its subdirectories
// are not checked.
func FindEnvs(path string, opts FindOpts) ([]*v1alpha1.Environment, error) {
	return findEnvsFromPaths([]string{path}, opts)
}

// FindEnvsFromPaths does the same as FindEnvs but takes a list of paths instead
func FindEnvsFromPaths(paths []string, opts FindOpts) ([]*v1alpha1.Environment, error) {
	return findEnvsFromPaths(paths, opts)
}

type findJsonnetFilesOut struct {
	jsonnetFiles []string
	err          error
}

type findEnvsOut struct {
	envs []*v1alpha1.Environment
	err  error
}

func findEnvsFromPaths(paths []string, opts FindOpts) ([]*v1alpha1.Environment, error) {
	if opts.Parallelism >= 0 {
		opts.Parallelism = runtime.NumCPU()
	}

	log.Debug().Int("parallelism", opts.Parallelism).Int("paths", len(paths)).Msg("Finding Tanka environments")
	startTime := time.Now()

	// find all jsonnet files within given paths
	pathChan := make(chan string, len(paths))
	findJsonnetFilesChan := make(chan findJsonnetFilesOut)
	for i := 0; i < opts.Parallelism; i++ {
		go func() {
			for path := range pathChan {
				jsonnetFiles, err := jsonnet.FindFiles(path, nil)
				var mainFiles []string
				for _, file := range jsonnetFiles {
					if filepath.Base(file) == jpath.DefaultEntrypoint {
						mainFiles = append(mainFiles, file)
					}
				}
				findJsonnetFilesChan <- findJsonnetFilesOut{jsonnetFiles: mainFiles, err: err}
			}
		}()
	}

	// push paths to channel
	var pathMap = map[string]bool{} // prevent duplicates
	for _, path := range paths {
		if _, ok := pathMap[path]; ok {
			continue
		}
		pathMap[path] = true
		pathChan <- path
	}

	// collect jsonnet files
	var jsonnetFiles []string
	for i := 0; i < len(paths); i++ {
		res := <-findJsonnetFilesChan
		if res.err != nil {
			return nil, res.err
		}
		jsonnetFiles = append(jsonnetFiles, res.jsonnetFiles...)
	}
	close(pathChan)
	close(findJsonnetFilesChan)

	findJsonnetFilesEndTime := time.Now()

	// find all environments within jsonnet files
	jsonnetFilesChan := make(chan string, len(jsonnetFiles))
	findEnvsChan := make(chan findEnvsOut)

	for i := 0; i < opts.Parallelism; i++ {
		go func() {
			for jsonnetFile := range jsonnetFilesChan {
				// try if this has envs
				list, err := List(jsonnetFile, Opts{JsonnetOpts: opts.JsonnetOpts})
				if err != nil &&
					// expected when looking for environments
					!errors.As(err, &jpath.ErrorNoBase{}) &&
					!errors.As(err, &jpath.ErrorFileNotFound{}) {
					findEnvsChan <- findEnvsOut{err: fmt.Errorf("%s:\n %w", jsonnetFile, err)}
					continue
				}
				filtered := []*v1alpha1.Environment{}
				// optionally filter
				if opts.Selector != nil && !opts.Selector.Empty() {
					for _, e := range list {
						if !opts.Selector.Matches(e.Metadata) {
							continue
						}
						filtered = append(filtered, e)
					}
				} else {
					filtered = append(filtered, list...)
				}
				findEnvsChan <- findEnvsOut{envs: filtered, err: nil}
			}
		}()
	}

	// push jsonnet files to channel
	var jsonnetFileMap = map[string]bool{} // prevent duplicates
	for _, jsonnetFile := range jsonnetFiles {
		if _, ok := jsonnetFileMap[jsonnetFile]; ok {
			continue
		}
		jsonnetFileMap[jsonnetFile] = true
		jsonnetFilesChan <- jsonnetFile
	}

	// collect environments
	var envs []*v1alpha1.Environment
	var errs []error
	for i := 0; i < len(jsonnetFiles); i++ {
		res := <-findEnvsChan
		if res.err != nil {
			errs = append(errs, res.err)
		}
		envs = append(envs, res.envs...)
	}
	close(jsonnetFilesChan)
	close(findEnvsChan)

	if len(errs) != 0 {
		return envs, ErrParallel{errors: errs}
	}

	findEnvsEndTime := time.Now()

	log.Info().
		Int("environments", len(envs)).
		Dur("ms_to_find_jsonnet_files", findJsonnetFilesEndTime.Sub(startTime)).
		Dur("ms_to_find_environments", findEnvsEndTime.Sub(findJsonnetFilesEndTime)).
		Msg("Found Tanka environments")

	return envs, nil
}
