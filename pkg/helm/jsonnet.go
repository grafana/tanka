package helm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// DefaultNameFormat to use when no nameFormat is supplied
const DefaultNameFormat = `{{ print .kind "_" .metadata.name | snakecase }}`

// JsonnetOpts are additional properties the consumer of the native func might
// pass.
type JsonnetOpts struct {
	TemplateOpts

	// CalledFrom is the file that calls helmTemplate. This is used to find the
	// vendored chart relative to this file
	CalledFrom string `json:"calledFrom"`
	// NameTemplate is used to create the keys in the resulting map
	NameFormat string `json:"nameFormat"`
}

// NativeFunc returns a jsonnet native function that provides the same
// functionality as `Helm.Template` of this package. Charts are required to be
// present on the local filesystem, at a relative location to the file that
// calls `helm.template()` / `std.native('helmTemplate')`. This guarantees
// hermeticity
func NativeFunc(h Helm) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: "helmTemplate",
		// Similar to `helm template [NAME] [CHART] [flags]` except 'conf' is a
		// bit more elaborate and chart is a local path
		Params: ast.Identifiers{"name", "chart", "opts"},
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
			opts, err := parseOpts(data[2])
			if err != nil {
				return "", err
			}

			// resolve the Chart relative to the caller
			callerDir := filepath.Dir(opts.CalledFrom)
			chart := filepath.Join(callerDir, chartpath)
			if _, err := os.Stat(chart); err != nil {
				return nil, fmt.Errorf("helmTemplate: Failed to find a Chart at '%s': %s. See https://tanka.dev/helm#failed-to-find-chart", chart, err)
			}

			// render resources
			list, err := h.Template(name, chart, opts.TemplateOpts)
			if err != nil {
				return nil, err
			}

			// convert list to map
			out, err := listAsMap(list, opts.NameFormat)
			if err != nil {
				return nil, err
			}

			return out, nil
		},
	}
}

func parseOpts(data interface{}) (*JsonnetOpts, error) {
	c, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var opts JsonnetOpts
	if err := json.Unmarshal(c, &opts); err != nil {
		return nil, err
	}

	// Charts are only allowed at relative paths. Use conf.CalledFrom to find the callers directory
	if opts.CalledFrom == "" {
		return nil, fmt.Errorf("helmTemplate: 'opts.calledFrom' is unset or empty.\nTanka needs this to find your Charts. See https://tanka.dev/helm#optscalledfrom-unset\n")
	}

	return &opts, nil
}

func listAsMap(list manifest.List, nameFormat string) (map[string]interface{}, error) {
	if nameFormat == "" {
		nameFormat = DefaultNameFormat
	}

	tmpl, err := template.New("").
		Funcs(sprig.TxtFuncMap()).
		Parse(nameFormat)
	if err != nil {
		return nil, fmt.Errorf("Parsing name format: %w", err)
	}

	out := make(map[string]interface{})
	for _, m := range list {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, m); err != nil {
			return nil, err
		}
		name := buf.String()

		if _, ok := out[name]; ok {
			return nil, ErrorDuplicateName{name: name, format: nameFormat}
		}
		out[name] = map[string]interface{}(m)
	}

	return out, nil
}

// ErrorDuplicateName means two resources share the same name using the given
// nameFormat.
type ErrorDuplicateName struct {
	name   string
	format string
}

func (e ErrorDuplicateName) Error() string {
	return fmt.Sprintf("Two resources share the same name '%s'. Please adapt the name template '%s'. See https://tanka.dev/helm#two-resources-share-the-same-name", e.name, e.format)
}
