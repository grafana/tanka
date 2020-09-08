package helm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

// TemplateOpts defines additional parameters that can be passed to the
// Helm.Template action
// TODO: Isolate between Helm.Template and NativeFunc
type TemplateOpts struct {
	// general
	Values map[string]interface{} `json:"values"`
	Flags  []string               `json:"flags"`

	// native func related
	CalledFrom string `json:"calledFrom"`
}

func confToArgs(conf TemplateOpts) ([]string, []string, error) {
	var args []string
	var tempFiles []string

	// create file and append to args
	if len(conf.Values) != 0 {
		valuesYaml, err := yaml.Marshal(conf.Values)
		if err != nil {
			return nil, nil, err
		}
		tmpFile, err := ioutil.TempFile(os.TempDir(), "tanka-")
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create temporary values.yaml")
		}
		tempFiles = append(tempFiles, tmpFile.Name())
		if _, err = tmpFile.Write(valuesYaml); err != nil {
			return nil, tempFiles, errors.Wrap(err, "failed to write to temporary values.yaml")
		}
		if err := tmpFile.Close(); err != nil {
			return nil, tempFiles, err
		}
		args = append(args, fmt.Sprintf("--values=%s", tmpFile.Name()))
	}

	// append custom flags to args
	args = append(args, conf.Flags...)

	return args, tempFiles, nil
}

// Template expands a Helm Chart into a regular manifest.List using the `helm
// template` command
func (h ExecHelm) Template(name, chart string, opts TemplateOpts) (manifest.List, error) {
	confArgs, tmpFiles, err := confToArgs(opts)
	if err != nil {
		return nil, err
	}
	for _, f := range tmpFiles {
		defer os.Remove(f)
	}

	args := []string{name, chart}
	args = append(args, confArgs...)

	cmd := h.cmd("template", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "Expanding Helm Chart")
	}

	var list manifest.List
	d := yaml.NewDecoder(&buf)
	for {
		var m manifest.Manifest
		if err := d.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "Parsing Helm output")
		}
		list = append(list, m)
	}

	return list, nil
}

// NativeFunc returns a jsonnet native function that provides the same
// functionality as `Helm.Template` of this package. Charts are required to be
// present on the local filesystem, at a relative location to the file that
// calls `helm.template()` / `std.native('helmTemplate')`. This guarantees
// hermeticity
func NativeFunc() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: "helmTemplate",
		// Similar to `helm template [NAME] [CHART] [flags]` except 'conf' is a
		// bit more elaborate and chart is a local path
		Params: ast.Identifiers{"name", "chart", "conf"},
		Func: func(data []interface{}) (interface{}, error) {
			name, ok := data[0].(string)
			if !ok {
				return nil, fmt.Errorf("First argument 'name' must be of 'string' type, got '%T' instead", data[0])
			}

			chartpath, ok := data[1].(string)
			if !ok {
				return nil, fmt.Errorf("Second argument 'chart' must be of 'string' type, got '%T' instead", data[1])
			}

			// TODO: validate data[2] actually follows the struct scheme
			c, err := json.Marshal(data[2])
			if err != nil {
				return "", err
			}
			var conf TemplateOpts
			if err := json.Unmarshal(c, &conf); err != nil {
				return "", err
			}

			// Charts are only allowed at relative paths. Use conf.CalledFrom to find the callers directory
			if conf.CalledFrom == "" {
				// TODO: rephrase and move lengthy explanation to website
				return nil, fmt.Errorf("helmTemplate: 'conf.calledFrom' is unset or empty.\nHowever, Tanka must know where helmTemplate was called from, to resolve the Helm Chart relative to that.\n")
			}
			callerDir := filepath.Dir(conf.CalledFrom)

			// resolve the Chart relative to the caller
			chart := filepath.Join(callerDir, chartpath)
			if _, err := os.Stat(chart); err != nil {
				// TODO: add website link for explanation
				return nil, fmt.Errorf("helmTemplate: Failed to find a Chart at '%s': %s", chart, err)
			}

			// TODO: Define Template on the Helm interface instead
			var h ExecHelm
			list, err := h.Template(name, chart, conf)
			if err != nil {
				return nil, err
			}

			out := make(map[string]interface{})
			for _, m := range list {
				name := fmt.Sprintf("%s_%s", m.Kind(), m.Metadata().Name())
				name = normalizeName(name)

				out[name] = map[string]interface{}(m)
			}

			return out, nil
		},
	}
}

func normalizeName(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ToLower(s)
	return s
}
