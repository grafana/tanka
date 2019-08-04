package jsonnet

import (
	"io/ioutil"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/sh0rez/tanka/pkg/native"
)

// EvaluateFile opens the file, reads it into memory and evaluates it afterwards (`Evaluate()`)
func EvaluateFile(jsonnetFile string) (string, error) {
	bytes, err := ioutil.ReadFile(jsonnetFile)
	if err != nil {
		return "", err
	}

	filename := filepath.Base(jsonnetFile)
	jpath, _, _ := jpath.Resolve(filepath.Dir(jsonnetFile), filename)

	return Evaluate(string(bytes), filename, jpath)
}

// Evaluate renders the given jssonet into a string
func Evaluate(sonnet string, filename string, jpath []string) (string, error) {
	importer := jsonnet.FileImporter{
		JPaths: jpath,
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&importer)
	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	return vm.EvaluateSnippet(filename, sonnet)
}
