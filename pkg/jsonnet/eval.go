package jsonnet

import (
	"io/ioutil"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/jsonnet/native"
)

// Modifiers allow to set optional paramters on the Jsonnet VM.
// See jsonnet.With* for this.
type Modifier func(vm *jsonnet.VM) error

// EvaluateFile opens the file, reads it into memory and evaluates it afterwards (`Evaluate()`)
func EvaluateFile(jsonnetFile string, mods ...Modifier) (string, error) {
	bytes, err := ioutil.ReadFile(jsonnetFile)
	if err != nil {
		return "", err
	}

	jpath, _, _, err := jpath.Resolve(filepath.Dir(jsonnetFile))
	if err != nil {
		return "", errors.Wrap(err, "resolving jpath")
	}
	return Evaluate(string(bytes), jpath, mods...)
}

// Evaluate renders the given jsonnet into a string
func Evaluate(sonnet string, jpath []string, mods ...Modifier) (string, error) {
	vm := jsonnet.MakeVM()
	vm.Importer(NewExtendedImporter(jpath))

	for _, mod := range mods {
		if err := mod(vm); err != nil {
			return "", err
		}
	}

	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	return vm.EvaluateSnippet("main.jsonnet", sonnet)
}

// WithExtCode allows to make the supplied snippet available to Jsonnet as an
// ext var
func WithExtCode(key, code string) Modifier {
	return func(vm *jsonnet.VM) error {
		vm.ExtCode(key, code)
		return nil
	}
}
