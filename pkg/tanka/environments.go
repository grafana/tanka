package tanka

import (
	"os"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

const BASEDIR_INDICATOR = "main.jsonnet"

// FindBaseDirs searches for possible environments
func FindBaseDirs(workdir string) (dirs []string, err error) {
	_, _, _, err = jpath.Resolve(workdir)
	if err == jpath.ErrorNoRoot {
		return nil, err
	}

	if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(filepath.Join(path, BASEDIR_INDICATOR)); err != nil {
			// missing file, not a valid environment directory
			return nil
		}
		dirs = append(dirs, path)
		return nil
	}); err != nil {
		return nil, err
	}
	return dirs, nil
}
