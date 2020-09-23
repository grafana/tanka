package jsonnet

import (
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
	ExtCode     InjectedCode
	TLACode     InjectedCode
	ImportPaths []string
	EvalPattern string
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

	return vm
}

// EvaluateFile evaluates the Jsonnet code in the given file and returns the
// result in JSON form. It disregards opts.ImportPaths in favor of automatically
// resolving these according to the specified file.
func EvaluateFile(jsonnetFile string, opts Opts) (string, error) {
	jpath, _, _, err := jpath.Resolve(jsonnetFile)
	if err != nil {
		return "", errors.Wrap(err, "resolving import paths")
	}
	opts.ImportPaths = jpath

	vm := MakeVM(opts)
	return vm.EvaluateFile(jsonnetFile)
}

// Evaluate renders the given jsonnet into a string
func Evaluate(filename, data string, opts Opts) (string, error) {
	jpath, _, _, err := jpath.Resolve(filename)
	if err != nil {
		return "", errors.Wrap(err, "resolving import paths")
	}
	opts.ImportPaths = jpath
	vm := MakeVM(opts)
	return vm.EvaluateAnonymousSnippet(filename, data)
}
