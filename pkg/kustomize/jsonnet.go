package kustomize

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
	// CalledFrom is the file that calls kustomizeBuild. This is used to find the
	// vendored Kustomize relative to this file
	CalledFrom string `json:"calledFrom"`
	// NameBuild is used to create the keys in the resulting map
	NameFormat string `json:"nameFormat"`
}

// NativeFunc returns a jsonnet native function that provides the same
// functionality as `Kustomize.Build` of this package. Kustomize yamls are required to be
// present on the local filesystem, at a relative location to the file that
// calls `kustomize.build()` / `std.native('kustomizeBuild')`. This guarantees
// hermeticity
func NativeFunc(k Kustomize) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: "kustomizeBuild",
		// Similar to `kustomize build {path}` where {path} is a local path
		Params: ast.Identifiers{"path", "opts"},
		Func: func(data []interface{}) (interface{}, error) {
			path, ok := data[0].(string)
			if !ok {
				return nil, fmt.Errorf("Argument 'path' must be of 'string' type, got '%T' instead", data[0])
			}

			// TODO: validate data[1] actually follows the struct scheme
			opts, err := parseOpts(data[1])
			if err != nil {
				return "", err
			}

			// resolve the Kustomize path relative to the caller
			callerDir := filepath.Dir(opts.CalledFrom)
			actual_path := filepath.Join(callerDir, path)
			if _, err := os.Stat(actual_path); err != nil {
				return nil, fmt.Errorf("kustomizeBuild: Failed to find kustomize at '%s': %s.", actual_path, err)
			}

			// render resources
			list, err := k.Build(actual_path)
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

	// Kustomize paths are only allowed at relative paths. Use conf.CalledFrom to find the callers directory
	if opts.CalledFrom == "" {
		return nil, fmt.Errorf("kustomizeBuild: 'opts.calledFrom' is unset or empty.\nTanka needs this to find your Kustomize.\n")
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
	return fmt.Sprintf("Two resources share the same name '%s'. Please adapt the name template '%s'.", e.name, e.format)
}
