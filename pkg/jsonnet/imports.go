package jsonnet

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/google/go-jsonnet/toolutils"
	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/jsonnet/native"
)

// TransitiveImports returns all recursive imports of an environment
func TransitiveImports(dir string) ([]string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	mainFile := filepath.Join(dir, "main.jsonnet")

	sonnet, err := ioutil.ReadFile(mainFile)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}

	jpath, _, rootDir, err := jpath.Resolve(dir)
	if err != nil {
		return nil, errors.Wrap(err, "resolving JPATH")
	}

	vm := jsonnet.MakeVM()
	vm.Importer(NewExtendedImporter(jpath))
	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	node, err := jsonnet.SnippetToAST("main.jsonnet", string(sonnet))
	if err != nil {
		return nil, errors.Wrap(err, "creating Jsonnet AST")
	}

	imports := make([]string, 0)
	if err = importRecursive(&imports, vm, node, "main.jsonnet"); err != nil {
		return nil, err
	}

	uniq := append(uniqueStringSlice(imports), mainFile)
	for i := range uniq {
		uniq[i], _ = filepath.Rel(rootDir, uniq[i])
	}
	sort.Strings(uniq)

	return uniq, nil
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

		abs, _ := filepath.Abs(foundAt)
		*list = append(*list, abs)

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

		abs, _ := filepath.Abs(foundAt)
		*list = append(*list, abs)

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
