package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

// findBaseDirs searches for possible environments
func findBaseDirs() (dirs []string) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	_, _, _, err = jpath.Resolve(pwd)
	if err == jpath.ErrorNoRoot {
		return
	}

	if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		requiredFiles := []string{"main.jsonnet", "spec.json"}
		for _, name := range requiredFiles {
			if _, err := os.Stat(filepath.Join(path, name)); err != nil {
				// missing file, not a valid environment directory
				return nil
			}
		}
		dirs = append(dirs, path)
		return nil
	}); err != nil {
		log.Fatalln(err)
	}
	return dirs
}
