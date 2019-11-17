package jsonnet

import (
	"encoding/json"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ExtendedImporter wraps jsonnet.FileImporter to add additional functionality:
// - `import "file.yaml"`
type ExtendedImporter struct {
	fi *jsonnet.FileImporter
}

// NewExtendedImporter returns a new instance of ExtendedImporter with the
// correct jpaths set up
func NewExtendedImporter(jpath []string) *ExtendedImporter {
	return &ExtendedImporter{
		fi: &jsonnet.FileImporter{
			JPaths: jpath,
		},
	}
}

func (i *ExtendedImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	// regularly import
	contents, foundAt, err = i.fi.Import(importedFrom, importedPath)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}

	// if yaml -> convert to json
	ext := filepath.Ext(foundAt)
	if ext == ".yaml" || ext == ".yml" {
		var data map[string]interface{}
		if err := yaml.Unmarshal([]byte(contents.String()), &data); err != nil {
			return jsonnet.Contents{}, "", errors.Wrapf(err, "unmarshalling yaml import '%s'", foundAt)
		}
		out, err := json.Marshal(data)
		if err != nil {
			return jsonnet.Contents{}, "", errors.Wrapf(err, "converting '%s' to json", foundAt)
		}
		contents = jsonnet.MakeContents(string(out))
	}

	return
}
