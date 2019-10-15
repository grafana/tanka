package jsonnet

import (
	"io/ioutil"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/google/go-jsonnet/toolutils"
	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/jsonnet/native"
)

// TransitiveImports returns all recursive imports of a file
func TransitiveImports(filename string) ([]string, error) {
	sonnet, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}

	jpath, _, _, err := jpath.Resolve(filepath.Dir(filename))
	if err != nil {
		return nil, errors.Wrap(err, "resolving JPATH")
	}
	importer := jsonnet.FileImporter{
		JPaths: jpath,
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&importer)
	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	node, err := jsonnet.SnippetToAST("main.jsonnet", string(sonnet))
	if err != nil {
		return nil, errors.Wrap(err, "creating Jsonnet AST")
	}

	imports := make([]string, 0)
	err = importRecursive(&imports, vm, node, "main.jsonnet")

	return uniqueStringSlice(imports), err
}

// importRecursive takes a Jsonnet VM and recursively imports the AST. Every
// found import is added to the `list` string slice, which will ultimately
// contain all recursive imports
func importRecursive(list *[]string, vm *jsonnet.VM, node ast.Node, currentPath string) error {
	switch node := node.(type) {
	// we have an `import`
	case *ast.Import:
		p := node.File.Value

		contents, foundAt, err := vm.ImportAST(currentPath, p)
		if err != nil {
			return errors.Wrap(err, "importing jsonnet")
		}

		*list = append(*list, foundAt)

		if err := importRecursive(list, vm, contents, foundAt); err != nil {
			return err
		}

	// we have an `importstr`
	case *ast.ImportStr:
		p := node.File.Value

		foundAt, err := vm.ResolveImport(currentPath, p)
		if err != nil {
			return errors.Wrap(err, "importing string")
		}

		*list = append(*list, foundAt)

	// neither `import` nor `importstr`, probably object or similar: try children
	default:
		for _, child := range toolutils.Children(node) {
			if err := importRecursive(list, vm, child, currentPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func uniqueStringSlice(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}
