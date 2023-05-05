package jsonnet

import (
	"os"
	"regexp"
	"time"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/grafana/tanka/pkg/jsonnet/implementation"
	"github.com/grafana/tanka/pkg/jsonnet/implementation/goimpl"
	"github.com/grafana/tanka/pkg/jsonnet/implementation/types"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

// Modifier allows to set optional parameters on the Jsonnet VM.
// See jsonnet.With* for this.
type Modifier func(vm *jsonnet.VM) error

// InjectedCode holds data that is "late-bound" into the VM
type InjectedCode map[string]string

// Set allows to set values on an InjectedCode, even when it is nil
func (i *InjectedCode) Set(key, value string) {
	if *i == nil {
		*i = make(InjectedCode)
	}

	(*i)[key] = value
}

// Opts are additional properties for the Jsonnet VM
type Opts struct {
	JsonnetImplementation string
	MaxStack              int
	ExtCode               InjectedCode
	TLACode               InjectedCode
	ImportPaths           []string
	EvalScript            string
	CachePath             string

	CachePathRegexes []*regexp.Regexp
}

// PathIsCached determines if a given path is matched by any of the configured cached path regexes
// If no path regexes are defined, all paths are matched
func (o Opts) PathIsCached(path string) bool {
	for _, regex := range o.CachePathRegexes {
		if regex.MatchString(path) {
			return true
		}
	}
	return len(o.CachePathRegexes) == 0
}

// Clone returns a deep copy of Opts
func (o Opts) Clone() Opts {
	extCode, tlaCode := InjectedCode{}, InjectedCode{}

	for k, v := range o.ExtCode {
		extCode[k] = v
	}

	for k, v := range o.TLACode {
		tlaCode[k] = v
	}

	return Opts{
		TLACode:     tlaCode,
		ExtCode:     extCode,
		ImportPaths: append([]string{}, o.ImportPaths...),
		EvalScript:  o.EvalScript,

		CachePath:        o.CachePath,
		CachePathRegexes: o.CachePathRegexes,
	}
}

// EvaluateFile evaluates the Jsonnet code in the given file and returns the
// result in JSON form. It disregards opts.ImportPaths in favor of automatically
// resolving these according to the specified file.
func EvaluateFile(jsonnetFile string, opts Opts) (string, error) {
	evalFunc := func(vm types.JsonnetVM) (string, error) {
		return vm.EvaluateFile(jsonnetFile)
	}
	data, err := os.ReadFile(jsonnetFile)
	if err != nil {
		return "", err
	}
	return evaluateSnippet(evalFunc, jsonnetFile, string(data), opts)
}

// Evaluate renders the given jsonnet into a string
// If cache options are given, a hash from the data will be computed and
// the resulting string will be cached for future retrieval
func Evaluate(path, data string, opts Opts) (string, error) {
	evalFunc := func(vm types.JsonnetVM) (string, error) {
		return vm.EvaluateAnonymousSnippet(path, data)
	}
	return evaluateSnippet(evalFunc, path, data, opts)
}

type evalFunc func(vm types.JsonnetVM) (string, error)

func evaluateSnippet(evalFunc evalFunc, path, data string, opts Opts) (string, error) {
	var cache *FileEvalCache
	if opts.CachePath != "" && opts.PathIsCached(path) {
		cache = NewFileEvalCache(opts.CachePath)
	}

	// Create VM
	jsonnetImpl := implementation.Get(opts.JsonnetImplementation)
	jpath, _, _, err := jpath.Resolve(path, false)
	if err != nil {
		return "", errors.Wrap(err, "resolving import paths")
	}
	opts.ImportPaths = jpath
	vm := jsonnetImpl.MakeVM(opts.ImportPaths, opts.ExtCode, opts.TLACode, opts.MaxStack)
	importVM := goimpl.MakeRawVM(opts.ImportPaths, opts.ExtCode, opts.TLACode, opts.MaxStack) // TODO: use interface

	var hash string
	if cache != nil {
		startTime := time.Now()
		if hash, err = getSnippetHash(importVM, path, data); err != nil {
			return "", err
		}
		cacheLog := log.Debug().Str("path", path).Str("hash", hash).Dur("duration_ms", time.Since(startTime))
		if v, err := cache.Get(hash); err != nil {
			return "", err
		} else if v != "" {
			cacheLog.Bool("cache_hit", true).Msg("computed snippet hash")
			return v, nil
		}
		cacheLog.Bool("cache_hit", false).Msg("computed snippet hash")
	}

	content, err := evalFunc(vm)
	if err != nil {
		return "", err
	}

	if cache != nil {
		return content, cache.Store(hash, content)
	}

	return content, nil
}
