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

type ImportVisitor func(who, what string) error

func VisitImportsFile(jsonnetFile string, v ImportVisitor) error {
	bytes, err := ioutil.ReadFile(jsonnetFile)
	if err != nil {
		return err
	}

	jpath, _, _, err := jpath.Resolve(filepath.Dir(jsonnetFile))
	if err != nil {
		return err
	}
	return VisitImports(string(bytes), jpath, v)
}

func VisitImports(sonnet string, jpath []string, v ImportVisitor) error {
	importer := TraceImporter{
		JPaths:  jpath,
		Visitor: v,
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&importer)
	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	if _, err := vm.EvaluateSnippet("main.jsonnet", sonnet); err != nil {
		return err
	}
	return nil
}

type TraceImporter struct {
	JPaths   []string
	Visitor  ImportVisitor
	importer *jsonnet.FileImporter
}

func (t *TraceImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	if t.importer == nil {
		t.importer = &jsonnet.FileImporter{
			JPaths: t.JPaths,
		}
	}

	contents, foundAt, err = t.importer.Import(importedFrom, importedPath)
	if err := t.Visitor(importedFrom, foundAt); err != nil {
		return jsonnet.Contents{}, "", err
	}
	return
}
