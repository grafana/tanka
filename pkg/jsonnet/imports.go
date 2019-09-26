package jsonnet

import (
	"io/ioutil"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/google/go-jsonnet/parser"
	"github.com/grafana/tanka/pkg/jpath"
	"github.com/grafana/tanka/pkg/native"
)

// ImportVisitor is a function that is invoked on every found import and is
// passed `who` imported `what`.
type ImportVisitor func(who, what string) error

// VisitImportsFile wraps VisitImports to load the Jsonnet from a file
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

// VisitImports calls the ImportVisitor for every recursive import of the Jsonnet.
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

	node, err := jsonnet.SnippetToAST("main.jsonnet", sonnet)
	if err != nil {
		return err
	}

	return importRecursive(vm, node, "main.jsonnet")
}

// importRecursive takes a Jsonnet VM and recursively imports the AST.
// This is especially useful in combination with the TraceImporter below.
func importRecursive(vm *jsonnet.VM, node ast.Node, currentPath string) error {
	switch node := node.(type) {
	case *ast.Import:
		p := node.File.Value
		contents, foundAt, err := vm.ImportAST(currentPath, p)
		if err != nil {
			return err
		}

		if err := importRecursive(vm, contents, foundAt); err != nil {
			return err
		}
	case *ast.ImportStr:
		p := node.File.Value
		_, err := vm.ResolveImport(currentPath, p)
		if err != nil {
			return err
		}
	default:
		for _, child := range parser.Children(node) {
			if err := importRecursive(vm, child, currentPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// TraceImporter wraps a jsonnet.FileImporter but also records every import to
// the ImportVisitor.
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
