package helmraiser

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type TemplateOpts struct {
	Values Values
	Flags  []string
}

// Template uses `helm template` to expand a Helm Chart into a plain manifest.List
func (h Helm) Template(name, chart string, opts TemplateOpts) (manifest.List, error) {
	argv := []string{name, chart}
	argv = append(argv, opts.Flags...)

	// values.yml
	tmp, err := writeTmpFile("values.yml", []byte(opts.Values.String()))
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)
	argv = append(argv, "--values="+tmp)

	// construct command
	cmd := h.run("template", argv...)
	cmd.Stderr = os.Stderr

	var buf bytes.Buffer
	cmd.Stdout = &buf

	// run command
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "Expanding Helm Chart")
	}

	// parse output
	var list manifest.List
	if err := parseYAMLLikeJSON(buf.Bytes(), &list); err != nil {
		return nil, errors.Wrap(err, "Parsing Helm output")
	}

	return list, nil
}

// parseYAMLLikeJSON loads YAML, reformats to JSON and loads it again. What
// sounds like a useless waste of resources actually is very valuable, because
// the json parser uses a smaller set of possible types than the YAML parser
// does.
// Several packages, like `google/go-jsonnet` or our very own `manifest` require this.
func parseYAMLLikeJSON(data []byte, ptr interface{}) error {
	var tmp []interface{}
	d := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var m map[string]interface{}
		if err := d.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		tmp = append(tmp, m)
	}

	jsonData, err := json.Marshal(tmp)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, ptr)
}
