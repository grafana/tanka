package helm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

// Helm provides actions on Helm charts.
type Helm struct{}

func (h Helm) cmd(action string, args ...string) *exec.Cmd {
	argv := []string{action}
	argv = append(argv, args...)

	return helmCmd(argv...)
}

func helmCmd(args ...string) *exec.Cmd {
	binary := "helm"
	if env := os.Getenv("TANKA_HELM_PATH"); env != "" {
		binary = env
	}
	return exec.Command(binary, args...)
}

// TemplateOpts defines additional parameters that can be passed to the
// Helm.Template action
type TemplateOpts struct {
	Values map[string]interface{}
	Flags  []string
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
func (h Helm) Template(name, chart string, opts TemplateOpts) (manifest.List, error) {
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
// functionality as `Helm.Template` of this package
func NativeFunc() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: "helmTemplate",
		// Lines up with `helm template [NAME] [CHART] [flags]` except 'conf' is a bit more elaborate
		Params: ast.Identifiers{"name", "chart", "conf"},
		Func: func(data []interface{}) (interface{}, error) {
			name, ok := data[0].(string)
			if !ok {
				return nil, fmt.Errorf("First argument 'name' must be of 'string' type, got '%T' instead", data[0])
			}

			chart, ok := data[1].(string)
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

			var h Helm
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
