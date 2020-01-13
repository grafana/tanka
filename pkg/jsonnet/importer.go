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
	fi           *jsonnet.FileImporter
	interceptors []ImportInterceptor
	processors   []ImportProcessor
}

// ImportInterceptor are executed before the actual importing. If they return
// something, this value is used.
type ImportInterceptor func(importedFrom, importedPath string) (c *string, foundAt string, err error)

// ImportProcessor are executed after the file import and may modify the result
// further
type ImportProcessor func(contents, foundAt string) (c *string, err error)

// NewExtendedImporter returns a new instance of ExtendedImporter with the
// correct jpaths set up
func NewExtendedImporter(jpath []string) *ExtendedImporter {
	return &ExtendedImporter{
		fi: &jsonnet.FileImporter{
			JPaths: jpath,
		},
		interceptors: []ImportInterceptor{tkInterceptor},
		processors:   []ImportProcessor{yamlProcessor},
	}
}

// Import implements the functionality offered by the ExtendedImporter
func (i *ExtendedImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	// check if an interceptor handles this
	for _, interceptor := range i.interceptors {
		c, foundAt, err := interceptor(importedFrom, importedPath)
		switch {
		case err != nil:
			return jsonnet.Contents{}, "", err
		case c == nil:
			continue
		default:
			return jsonnet.MakeContents(*c), foundAt, nil
		}
	}

	// regularly import
	contents, foundAt, err = i.fi.Import(importedFrom, importedPath)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}

	// check if needs postprocessing
	for _, processor := range i.processors {
		c, err := processor(contents.String(), foundAt)
		switch {
		case err != nil:
			return jsonnet.Contents{}, "", err
		case c == nil:
			continue
		default:
			return jsonnet.MakeContents(*c), foundAt, nil
		}
	}

	return contents, foundAt, nil
}

// tkInterceptor provides `tk.libsonnet` from memory (builtin)
func tkInterceptor(importedFrom, importedPath string) (contents *string, foundAt string, err error) {
	if importedPath != "tk" {
		return nil, "", nil
	}

	s := tkLibsonnet
	return &s, filepath.Join(locationInternal, "tk.libsonnet"), nil
}

// yamlProcessor catches yaml files and converts them to JSON so that they can
// be used with `import`
func yamlProcessor(contents, foundAt string) (c *string, err error) {
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

	s := string(out)
	return &s, nil
}
