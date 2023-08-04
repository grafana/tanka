package jsonnet

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/jsonnet/native"
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
	MaxStack    int
	ExtCode     InjectedCode
	TLACode     InjectedCode
	ImportPaths []string
	EvalScript  string
	CachePath   string

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

type vmPool struct {
	mutex     sync.Mutex
	available map[string][]*jsonnet.VM
}

var VMPool = vmPool{
	available: map[string][]*jsonnet.VM{},
}

func optsHash(opts Opts) string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s", opts.ExtCode)))
	hash.Write([]byte(fmt.Sprintf("%s", opts.ImportPaths)))
	return hex.EncodeToString(hash.Sum(nil))
}

// Get returns a Jsonnet VM with some extensions of Tanka, including:
// - extended importer
// - extCode and tlaCode applied
// - native functions registered
// If a VM is available in the pool, it will be reused. Otherwise a new one will be created.
func (p *vmPool) Get(opts Opts) *jsonnet.VM {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	lastImport := opts.ImportPaths[len(opts.ImportPaths)-1]
	var vm *jsonnet.VM
	hash := optsHash(opts)
	if cached := p.available[hash]; len(cached) > 0 {
		log.Trace().Str("path", lastImport).Msg("reusing Jsonnet VM")
		vm = cached[0]
		p.available[hash] = cached[1:]
	} else {
		log.Trace().Str("path", lastImport).Msg("creating new Jsonnet VM")
		vm = jsonnet.MakeVM()
		if opts.MaxStack > 0 {
			vm.MaxStack = opts.MaxStack
		}
		for _, nf := range native.Funcs() {
			vm.NativeFunction(nf)
		}
		vm.Importer(NewExtendedImporter(opts.ImportPaths))

		for k, v := range opts.ExtCode {
			vm.ExtCode(k, v)
		}
	}

	for k, v := range opts.TLACode {
		vm.TLACode(k, v)
	}

	return vm
}

// Release returns a Jsonnet VM to the pool
func (p *vmPool) Release(vm *jsonnet.VM, opts Opts) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.available[optsHash(opts)] = append(p.available[optsHash(opts)], vm)
}

// EvaluateFile evaluates the Jsonnet code in the given file and returns the
// result in JSON form. It disregards opts.ImportPaths in favor of automatically
// resolving these according to the specified file.
func EvaluateFile(jsonnetFile string, opts Opts) (string, error) {
	evalFunc := func(vm *jsonnet.VM) (string, error) {
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
	evalFunc := func(vm *jsonnet.VM) (string, error) {
		return vm.EvaluateAnonymousSnippet(path, data)
	}
	return evaluateSnippet(evalFunc, path, data, opts)
}

type evalFunc func(vm *jsonnet.VM) (string, error)

func evaluateSnippet(evalFunc evalFunc, path, data string, opts Opts) (string, error) {
	var cache *FileEvalCache
	if opts.CachePath != "" && opts.PathIsCached(path) {
		cache = NewFileEvalCache(opts.CachePath)
	}

	// Create VM
	jpath, _, _, err := jpath.Resolve(path, false)
	if err != nil {
		return "", errors.Wrap(err, "resolving import paths")
	}
	opts.ImportPaths = jpath
	vm := VMPool.Get(opts)
	defer VMPool.Release(vm, opts)

	var hash string
	if cache != nil {
		startTime := time.Now()
		if hash, err = getSnippetHash(vm, path, data); err != nil {
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
