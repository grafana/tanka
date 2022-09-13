package jsonnet

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/jsonnet/native"
)

var importsRegexp = regexp.MustCompile(`import(str)?\s+['"]([^'"%()]+)['"]`)

// TransitiveImports returns all recursive imports of an environment
func TransitiveImports(dir string) ([]string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return nil, err
	}

	entrypoint, err := jpath.Entrypoint(dir)
	if err != nil {
		return nil, err
	}

	jpath, _, rootDir, err := jpath.Resolve(dir, false)
	if err != nil {
		return nil, errors.Wrap(err, "resolving JPATH")
	}

	vm := jsonnet.MakeVM()
	vm.Importer(NewExtendedImporter(jpath))
	for _, nf := range native.Funcs() {
		vm.NativeFunction(nf)
	}

	imports := make(map[string]bool)
	if err = importRecursiveStrict(imports, vm, entrypoint); err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(imports)+1)
	for k := range imports {

		// Try to resolve any symlinks; use the original path as a last resort
		p, err := filepath.EvalSymlinks(k)
		if err != nil {
			return nil, errors.Wrap(err, "resolving symlinks")
		}
		paths = append(paths, p)

	}
	paths = append(paths, entrypoint)

	for i := range paths {
		paths[i], _ = filepath.Rel(rootDir, paths[i])

		// Normalize path separators for windows
		paths[i] = filepath.ToSlash(paths[i])
	}
	sort.Strings(paths)

	return paths, nil
}

// importRecursiveStrict does the same as importRecursive, but returns an error
// if a file is not found during when importing
func importRecursiveStrict(list map[string]bool, vm *jsonnet.VM, filename string) error {
	return importRecursive(list, vm, filename, false)
}

// importRecursive takes a Jsonnet VM and recursively imports the AST. Every
// found import is added to the `list` string slice, which will ultimately
// contain all recursive imports
func importRecursive(list map[string]bool, vm *jsonnet.VM, filename string, ignoreMissing bool) error {

	content, err := os.ReadFile(filename)
	if err != nil {
		if res, err := os.Stat(filename); err == nil && res.IsDir() {
			filename, err = jpath.Entrypoint(filename)
			if err != nil {
				return err
			}
			content, err = os.ReadFile(filename)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("reading file %s: %s", filename, err)
		}
	}

	matches := importsRegexp.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		foundAt, err := vm.ResolveImport(filename, match[2])
		if err != nil {
			if ignoreMissing {
				continue
			}
			return err
		}

		abs, _ := filepath.Abs(foundAt)

		if list[abs] {
			return nil
		}
		list[abs] = true

		if match[1] == "str" {
			continue
		}

		if err := importRecursive(list, vm, abs, ignoreMissing); err != nil {
			return err
		}
	}
	return nil
}

var fileHashes sync.Map

// getSnippetHash takes a jsonnet snippet and calculates a hash from its content
// and the content of all of its dependencies.
// File hashes are cached in-memory to optimize multiple executions of this function in a process
func getSnippetHash(vm *jsonnet.VM, path, data string) (string, error) {
	result := map[string]bool{}
	if err := importRecursive(result, vm, path, true); err != nil {
		return "", err
	}
	fileNames := []string{}
	for file := range result {
		fileNames = append(fileNames, file)
	}
	sort.Strings(fileNames)

	fullHasher := sha256.New()
	fullHasher.Write([]byte(data))
	for _, file := range fileNames {
		var fileHash []byte
		if got, ok := fileHashes.Load(file); ok {
			fileHash = got.([]byte)
		} else {
			bytes, err := os.ReadFile(file)
			if err != nil {
				return "", err
			}
			hash := sha256.New()
			fileHash = hash.Sum(bytes)
			fileHashes.Store(file, fileHash)
		}
		fullHasher.Write(fileHash)
	}

	return base64.URLEncoding.EncodeToString(fullHasher.Sum(nil)), nil
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
