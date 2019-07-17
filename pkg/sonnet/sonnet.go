package sonnet

import (
	"io/ioutil"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/sh0rez/tanka/pkg/native"
)

func EvaluateFile(jsonnetFile string) (string, error) {
	bytes, err := ioutil.ReadFile(jsonnetFile)
	if err != nil {
		return "", err
	}

	jpath, _, _ := jpath.Resolve(filepath.Dir(jsonnetFile))
	return Evaluate(string(bytes), jpath)
}

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
