package jsonnet

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

// FindImporterForFiles finds the entrypoints (main.jsonnet files) that import the given files.
// It looks through imports transitively, so if a file is imported through a chain, it will still be reported.
// If the given file is a main.jsonnet file, it will be returned as well.
func FindImporterForFiles(root string, files []string, chain map[string]struct{}) ([]string, error) {
	if chain == nil {
		chain = make(map[string]struct{})
	}

	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	importers := map[string]struct{}{}

	if len(chain) == 0 {
		for i := range files {
			files[i], err = filepath.Abs(files[i])
			if err != nil {
				return nil, err
			}

			symlink, err := evalSymlinks(files[i])
			if err != nil {
				return nil, err
			}
			if symlink != files[i] {
				files = append(files, symlink)
			}

			symlinks, err := findSymlinks(root, files[i])
			if err != nil {
				return nil, err
			}
			files = append(files, symlinks...)
		}

		files = uniqueStringSlice(files)
	}

	for _, file := range files {
		if filepath.Base(file) == jpath.DefaultEntrypoint {
			importers[file] = struct{}{}
		}

		newImporters, err := findImporters(root, file, chain)
		if err != nil {
			return nil, err
		}
		for _, importer := range newImporters {
			importers[importer] = struct{}{}
		}
	}

	var importersSlice []string
	for importer := range importers {
		importersSlice = append(importersSlice, importer)
	}

	sort.Strings(importersSlice)

	return importersSlice, nil
}

type cachedJsonnetFile struct {
	Base       string
	Imports    []string
	Content    string
	IsMainFile bool
}

var jsonnetFilesMap = make(map[string]map[string]*cachedJsonnetFile)
var symlinkCache = make(map[string]string)

func evalSymlinks(path string) (string, error) {
	var err error
	eval, ok := symlinkCache[path]
	if !ok {
		eval, err = filepath.EvalSymlinks(path)
		if err != nil {
			return "", err
		}
		symlinkCache[path] = eval
	}
	return eval, nil
}

func findSymlinks(root, file string) ([]string, error) {
	var symlinks []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			eval, err := evalSymlinks(path)
			if err != nil {
				return err
			}
			if strings.Contains(file, eval) {
				symlinks = append(symlinks, strings.Replace(file, eval, path, 1))
			}
		}

		return nil
	})

	return symlinks, err
}

func findImporters(root string, searchForFile string, chain map[string]struct{}) ([]string, error) {
	if _, ok := chain[searchForFile]; ok {
		return nil, nil
	}
	chain[searchForFile] = struct{}{}

	if _, ok := jsonnetFilesMap[root]; !ok {
		jsonnetFilesMap[root] = make(map[string]*cachedJsonnetFile)

		files, err := FindFiles(root, nil)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}
			matches := importsRegexp.FindAllStringSubmatch(string(content), -1)

			cachedObj := &cachedJsonnetFile{
				Content:    string(content),
				IsMainFile: strings.HasSuffix(file, jpath.DefaultEntrypoint),
			}
			for _, match := range matches {
				cachedObj.Imports = append(cachedObj.Imports, match[2])
			}
			jsonnetFilesMap[root][file] = cachedObj
		}
	}
	jsonnetFiles := jsonnetFilesMap[root]

	var importers []string
	var intermediateImporters []string

	for jsonnetFilePath, jsonnetFileContent := range jsonnetFiles {
		isImporter := false
		for _, importPath := range jsonnetFileContent.Imports {
			if filepath.Base(importPath) != filepath.Base(searchForFile) { // If the filename is not the same as the file we are looking for, skip
				continue
			}

			// Match on relative imports with ..
			// Jsonnet also matches all intermediary paths for some reason, so we look at them too
			doubleDotCount := strings.Count(importPath, "..")
			if doubleDotCount > 0 {
				importPath = strings.ReplaceAll(importPath, "../", "")
				for i := 0; i <= doubleDotCount; i++ {
					dir := filepath.Dir(jsonnetFilePath)
					for j := 0; j < i; j++ {
						dir = filepath.Dir(dir)
					}
					testImportPath := filepath.Join(dir, importPath)
					isImporter = pathMatches(searchForFile, testImportPath)
				}
			}

			// Match on imports to lib/ or vendor/
			if !isImporter {
				importPath = strings.ReplaceAll(importPath, "./", "")
				isImporter = pathMatches(searchForFile, filepath.Join(root, "vendor", importPath)) || pathMatches(searchForFile, filepath.Join(root, "lib", importPath))
			}

			// Match on imports to the base dir where the file is located (e.g. in the env dir)
			if !isImporter {
				if jsonnetFileContent.Base == "" {
					base, err := jpath.FindBase(jsonnetFilePath, root)
					if err != nil {
						return nil, err
					}
					jsonnetFileContent.Base = base
				}
				isImporter = strings.HasPrefix(searchForFile, jsonnetFileContent.Base) && strings.HasSuffix(searchForFile, importPath)
			}

			if isImporter {
				if jsonnetFileContent.IsMainFile {
					importers = append(importers, jsonnetFilePath)
				} else {
					intermediateImporters = append(intermediateImporters, jsonnetFilePath)
				}
				break
			}
		}
	}

	if len(intermediateImporters) > 0 {
		newImporters, err := FindImporterForFiles(root, intermediateImporters, chain)
		if err != nil {
			return nil, err
		}
		importers = append(importers, newImporters...)
	}

	return importers, nil
}

func pathMatches(path1, path2 string) bool {
	if path1 == path2 {
		return true
	}

	var err error

	evalPath1, err := evalSymlinks(path1)
	if err != nil {
		return false
	}

	evalPath2, err := evalSymlinks(path2)
	if err != nil {
		return false
	}

	return evalPath1 == evalPath2
}
