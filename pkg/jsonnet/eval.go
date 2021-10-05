package jsonnet

import (
	"io/ioutil"
	"regexp"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"

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
func (opts Opts) PathIsCached(path string) bool {
	for _, regex := range opts.CachePathRegexes {
		if regex.MatchString(path) {
			return true
		}
	}
	return len(opts.CachePathRegexes) == 0
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

// MakeVM returns a Jsonnet VM with some extensions of Tanka, including:
// - extended importer
// - extCode and tlaCode applied
// - native functions registered
func MakeVM(opts Opts) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.Importer(NewExtendedImporter(opts.ImportPaths))

	for k, v := range opts.ExtCode {
		vm.ExtCode(k, v)
	}
	for k, v := range opts.TLACode {
		vm.TLACode(k, v)
	}

	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	if opts.MaxStack > 0 {
		vm.MaxStack = opts.MaxStack
	}

	return vm
}

// EvaluateFile evaluates the Jsonnet code in the given file and returns the
// result in JSON form. It disregards opts.ImportPaths in favor of automatically
// resolving these according to the specified file.
func EvaluateFile(jsonnetFile string, opts Opts) (string, error) {
	bytes, _ := ioutil.ReadFile(jsonnetFile)
	return Evaluate(jsonnetFile, string(bytes), opts)

}

// Evaluate renders the given jsonnet into a string
// If cache options are given, a hash from the data will be computed and
//  the resulting string will be cached for future retrieval
func Evaluate(path, data string, opts Opts) (string, error) {
	var cache *FileEvalCache
	if opts.CachePath != "" && opts.PathIsCached(path) {
		cache = NewFileEvalCache(opts.CachePath)
	}

	// Create VM
	jpath, _, _, err := jpath.Resolve(path)
	if err != nil {
		return "", errors.Wrap(err, "resolving import paths")
	}
	opts.ImportPaths = jpath
	vm := MakeVM(opts)

	var hash string
	if cache != nil {
		if hash, err = getSnippetHash(vm, path, data); err != nil {
			return "", err
		}
		if v, err := cache.Get(hash); err != nil {
			return "", err
		} else if v != "" {
			return v, nil
		}
	}

	content, err := vm.EvaluateAnonymousSnippet(path, data)
	if err != nil {
		return "", err
	}

	if cache != nil {
		return content, cache.Store(hash, content)
	}

	return content, nil
}
