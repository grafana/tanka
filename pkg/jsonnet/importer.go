package jsonnet

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const locationInternal = "<internal>"

// ExtendedImporter wraps jsonnet.FileImporter to add additional functionality:
// - `import "file.yaml"`
// - `import "tk"`
type ExtendedImporter struct {
	loaders    []importLoader    // for loading jsonnet from somewhere. First one that returns non-nil is used
	processors []importProcessor // for post-processing (e.g. yaml -> json)
}

// importLoader are executed before the actual importing. If they return
// something, this value is used.
type importLoader func(importedFrom, importedPath string) (c *jsonnet.Contents, foundAt string, err error)

// importProcessor are executed after the file import and may modify the result
// further
type importProcessor func(contents, foundAt string) (c *jsonnet.Contents, err error)

// NewExtendedImporter returns a new instance of ExtendedImporter with the
// correct jpaths set up
func NewExtendedImporter(jpath []string) *ExtendedImporter {
	return &ExtendedImporter{
		loaders: []importLoader{
			tkLoader,
			newFileLoader(&jsonnet.FileImporter{
				JPaths: jpath,
			})},
		processors: []importProcessor{
			// TODO: re-enable this once we can without side-effects
			// (https://github.com/grafana/tanka/issues/135)
			//
			// yamlProcessor,
		},
	}
}

// Import implements the functionality offered by the ExtendedImporter
func (i *ExtendedImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	// load using loader
	for _, loader := range i.loaders {
		c, f, err := loader(importedFrom, importedPath)
		if err != nil {
			return jsonnet.Contents{}, "", err
		}
		if c != nil {
			contents = *c
			foundAt = f
			break
		}
	}

	// check if needs postprocessing
	for _, processor := range i.processors {
		c, err := processor(contents.String(), foundAt)
		if err != nil {
			return jsonnet.Contents{}, "", err
		}
		if c != nil {
			contents = *c
			break
		}
	}

	return contents, foundAt, nil
}

// tkLoader provides `tk.libsonnet` from memory (builtin)
func tkLoader(importedFrom, importedPath string) (contents *jsonnet.Contents, foundAt string, err error) {
	if importedPath != "tk" {
		return nil, "", nil
	}

	s := jsonnet.MakeContents(tkLibsonnet)
	return &s, filepath.Join(locationInternal, "tk.libsonnet"), nil
}

// newFileLoader returns an importLoader that uses jsonnet.FileImporter to source
// files from the local filesystem
func newFileLoader(fi *jsonnet.FileImporter) importLoader {
	return func(importedFrom, importedPath string) (contents *jsonnet.Contents, foundAt string, err error) {
		var c jsonnet.Contents
		c, foundAt, err = fi.Import(importedFrom, importedPath)
		return &c, foundAt, err
	}
}

// yamlProcessor catches yaml files and converts them to JSON so that they can
// be used with `import`
func yamlProcessor(contents, foundAt string) (c *jsonnet.Contents, err error) {
	ext := filepath.Ext(foundAt)
	if yaml := ext == ".yaml" || ext == ".yml"; !yaml {
		return nil, nil
	}

	ret := []interface{}{}
	d := yaml.NewDecoder(strings.NewReader(contents))
	for {
		var doc interface{}
		if err := d.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrapf(err, "unmarshalling yaml import '%s'", foundAt)
		}
		ret = append(ret, doc)
	}

	var data interface{} = ret
	if len(ret) == 1 {
		data = ret[0]
	}

	out, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrapf(err, "converting '%s' to json", foundAt)
	}

	s := jsonnet.MakeContents(string(out))
	return &s, nil
}
