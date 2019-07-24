package jsonnet

import (
	"io/ioutil"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"
	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/sh0rez/tanka/pkg/native"
)

// EvaluateFile opens the file, reads it into memory and evaluates it afterwards (`Evaluate()`)
func EvaluateFile(jsonnetFile string) (string, error) {
	bytes, err := ioutil.ReadFile(jsonnetFile)
	if err != nil {
		return "", err
	}

	jpath, _, _, err := jpath.Resolve(filepath.Dir(jsonnetFile))
	if err != nil {
		return "", errors.Wrap(err, "resolving jpath")
	}
	return Evaluate(string(bytes), jpath)
}

// Evaluate renders the given jssonet into a string
func Evaluate(sonnet string, jpath []string) (string, error) {
	importer := jsonnet.FileImporter{
		JPaths: jpath,
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&importer)
	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	return vm.EvaluateSnippet("main.jsonnet", sonnet)
}
