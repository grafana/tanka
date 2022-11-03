package jsonnet

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

var (
	importersCache    = make(map[string][]string)
	jsonnetFilesCache = make(map[string]map[string]*cachedJsonnetFile)
	symlinkCache      = make(map[string]string)
)

type cachedJsonnetFile struct {
	Base       string
	Imports    []string
	Content    string
	IsMainFile bool
}

// FindImporterForFiles finds the entrypoints (main.jsonnet files) that import the given files.
// It looks through imports transitively, so if a file is imported through a chain, it will still be reported.
// If the given file is a main.jsonnet file, it will be returned as well.
func FindImporterForFiles(root string, files []string) ([]string, error) {
	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	importers := map[string]struct{}{}

	if files, err = expandSymlinksInFiles(root, files); err != nil {
		return nil, err
	}

	// Loop through all given files and add their importers to the list
	for _, file := range files {
		if filepath.Base(file) == jpath.DefaultEntrypoint {
			importers[file] = struct{}{}
		}

		newImporters, err := findImporters(root, file, map[string]struct{}{})
		if err != nil {
			return nil, err
		}
		for _, importer := range newImporters {
			importers[importer] = struct{}{}
		}
	}

	return mapToArray(importers), nil
}

// expandSymlinksInFiles takes an array of files and adds to it:
// - all symlinks that point to the files
// - all files that are pointed to by the symlinks
func expandSymlinksInFiles(root string, files []string) ([]string, error) {
	filesMap := map[string]struct{}{}

	for _, file := range files {
		file, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}
		filesMap[file] = struct{}{}

		symlink, err := evalSymlinks(file)
		if err != nil {
			return nil, err
		}
		if symlink != file {
			filesMap[symlink] = struct{}{}
		}

		symlinks, err := findSymlinks(root, file)
		if err != nil {
			return nil, err
		}
		for _, symlink := range symlinks {
			filesMap[symlink] = struct{}{}
		}
	}

	return mapToArray(filesMap), nil
}

// evalSymlinks returns the path after following all symlinks.
// It caches the results to avoid unnecessary work.
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

// findSymlinks finds all symlinks that point to the given file.
// It's restricted to the given root directory.
// It's used in the case where a user wants to find which entrypoints import a given file.
// In that case, we also want to find the entrypoints that import a symlink to the file.
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
	// If we've already looked through this file in the current execution, don't do it again and return an empty list to end the recursion
	// Jsonnet supports cyclic imports (as long as the _attributes_ being used are not cyclic)
	if _, ok := chain[searchForFile]; ok {
		return nil, nil
	}
	chain[searchForFile] = struct{}{}

	// If we've already computed the importers for a file, return the cached result
	key := root + ":" + searchForFile
	if importers, ok := importersCache[key]; ok {
		return importers, nil
	}

	jsonnetFiles, err := createJsonnetFileCache(root)
	if err != nil {
		return nil, err
	}

	var importers []string
	var intermediateImporters []string

	for jsonnetFilePath, jsonnetFileContent := range jsonnetFiles {
		isImporter := false
		// For all imports in all jsonnet files, check if they import the file we're looking for
		for _, importPath := range jsonnetFileContent.Imports {
			// If the filename is not the same as the file we are looking for, skip it
			if filepath.Base(importPath) != filepath.Base(searchForFile) {
				continue
			}

			// Remove any `./` or `../` that can be removed just by looking at the given path
			// ex: `./foo/bar.jsonnet` -> `foo/bar.jsonnet` or `/foo/../bar.jsonnet` -> `/bar.jsonnet`
			importPath = filepath.Clean(importPath)

			// Match on relative imports with ..
			// Jsonnet also matches relative imports that are one level deeper than they should be
			// Example: Given two envs (env1 and env2), the two following imports in `env1/main.jsonnet will work`: `../env2/main.jsonnet` and `../../env2/main.jsonnet`
			// This can lead to false positives, but ruling them out would require much more complex logic
			if strings.HasPrefix(importPath, "..") {
				shallowImport := filepath.Clean(filepath.Join(filepath.Dir(jsonnetFilePath), strings.Replace(importPath, "../", "", 1)))
				importPath = filepath.Clean(filepath.Join(filepath.Dir(jsonnetFilePath), importPath))

				isImporter = pathMatches(searchForFile, importPath) || pathMatches(searchForFile, shallowImport)
			}

			// Match on imports to lib/ or vendor/
			if !isImporter {
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

			// If the file we're looking in imports one of the files we're looking for, add it to the list
			// It can either be an importer that we want to return (from a main file) or an intermediate importer
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

	// Process intermediate importers recursively
	// This will go on until we hit a main file, which will be returned
	if len(intermediateImporters) > 0 {
		for _, intermediateImporter := range intermediateImporters {
			newImporters, err := findImporters(root, intermediateImporter, chain)
			if err != nil {
				return nil, err
			}
			importers = append(importers, newImporters...)
		}
	}

	importersCache[key] = importers
	return importers, nil
}

func createJsonnetFileCache(root string) (map[string]*cachedJsonnetFile, error) {
	if val, ok := jsonnetFilesCache[root]; ok {
		return val, nil
	}
	jsonnetFilesCache[root] = make(map[string]*cachedJsonnetFile)

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
		jsonnetFilesCache[root][file] = cachedObj
	}

	return jsonnetFilesCache[root], nil
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

func mapToArray(m map[string]struct{}) []string {
	var arr []string
	for k := range m {
		arr = append(arr, k)
	}
	sort.Strings(arr)
	return arr
}
