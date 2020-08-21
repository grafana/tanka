package helmraiser

import (
	"encoding/json"
	"fmt"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// JsonnetOpts are the options accepted from Jsonnet
type JsonnetOpts struct {
	Values map[string]interface{} `json:"values"`
	Flags  []string               `json:"flags"`
	Repo   Repo                   `json:"repo"`
}

// helmTemplate wraps and runs `helm template`
// returns the generated manifests in a map
func HelmTemplate() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name: "helmTemplate",
		// Lines up with `helm template [NAME] [CHART] [flags]` except 'conf' is a bit more elaborate
		Params: ast.Identifiers{"name", "chart", "conf"},
		Func: func(data []interface{}) (interface{}, error) {
			// parse params received from Jsonnet
			name, chart := data[0].(string), data[1].(string)
			c, err := json.Marshal(data[2])
			if err != nil {
				return "", err
			}
			var conf JsonnetOpts
			if err := json.Unmarshal(c, &conf); err != nil {
				return "", err
			}

			// construct helm handler
			h, err := NewHelm(Repos{conf.Repo})
			if err != nil {
				return "", err
			}
			defer h.Close()

			// expand chart to yaml
			list, err := h.Template(name, chart, TemplateOpts{
				Flags:  conf.Flags,
				Values: conf.Values,
			})
			if err != nil {
				return nil, err
			}

			// transform it to map for easier patching from Jsonnet
			m := listToMSI(list)
			return m, nil
		},
	}
}

func listToMSI(list manifest.List) map[string]interface{} {
	out := make(map[string]interface{})

	// snake_case string
	normalizeName := func(s string) string {
		s = strings.ReplaceAll(s, "-", "_")
		s = strings.ReplaceAll(s, ":", "_")
		s = strings.ToLower(s)
		return s
	}

	for _, m := range list {
		name := fmt.Sprintf("%s_%s", m.Kind(), m.Metadata().Name())
		name = normalizeName(name)

		out[name] = map[string]interface{}(m)
	}

	return out
}
