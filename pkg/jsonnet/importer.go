package jsonnet

import (
	"bytes"
	"encoding/json"
	"io"
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

// Import implements the functionality offered by the ExtendedImporter
func (i *ExtendedImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	// regularly import
	contents, foundAt, err = i.fi.Import(importedFrom, importedPath)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}

	// if yaml -> convert to json
	ext := filepath.Ext(foundAt)
	if ext == ".yaml" || ext == ".yml" {
		ret := []interface{}{}
		d := yaml.NewDecoder(bytes.NewReader([]byte(contents.String())))
		for {
			var doc interface{}
			if err := d.Decode(&doc); err != nil {
				if err == io.EOF {
					break
				}
				return jsonnet.Contents{}, "", errors.Wrapf(err, "unmarshalling yaml import '%s'", foundAt)
			}
			ret = append(ret, doc)
		}
		var data interface{}
		if len(ret) == 1 {
			data = ret[0]
		} else {
			data = ret
		}
		out, err := json.Marshal(data)
		if err != nil {
			return jsonnet.Contents{}, "", errors.Wrapf(err, "converting '%s' to json", foundAt)
		}
		contents = jsonnet.MakeContents(string(out))
	}

	return
}
