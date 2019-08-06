package jsonnet

import (
	"path/filepath"
)

// TransitiveImports returns a slice with all files this file imports plus downstream imports
func TransitiveImports(filename string) ([]string, error) {
	imports := map[string][]string{}
	if err := VisitImportsFile(filename, func(who, what string) error {
		if imports[who] == nil {
			imports[who] = []string{}
		}
		imports[who] = append(imports[who], what)
		return nil
	}); err != nil {
		return nil, err
	}

	deps := map[string]*File{}
	for k, v := range imports {
		deps[k] = &File{Imports: v}
	}

	for _, d := range deps {
		resolveTransitives(d, deps)
	}

	for _, d := range deps {
		d.Dependencies = uniqueStringSlice(d.Dependencies)
	}

	return deps[filepath.Base(filename)].Dependencies, nil
}

// File represents a jsonnet file that may import other files
type File struct {
	// List of files this file imports
	Imports []string
	// Full list of transitive imports
	Dependencies []string
}

func resolveTransitives(f *File, deps map[string]*File) {
	// already resolved
	if len(f.Dependencies) != 0 {
		return
	}

	for _, i := range f.Imports {
		f.Dependencies = append(f.Dependencies, i)

		// import has no dependencies
		if deps[i] == nil {
			continue
		}

		// import dependencies have not yet been resolved
		if len(deps[i].Dependencies) == 0 {
			resolveTransitives(deps[i], deps)
		}

		f.Dependencies = append(f.Dependencies, deps[i].Dependencies...)
	}
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
