package kustomize

import (
	"bytes"
	"io"
	"os"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

// Build expands a Kustomize into a regular manifest.List using the `kustomize
// build` command
func (k ExecKustomize) Build(path string) (manifest.List, error) {
	cmd := k.cmd("build", path)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "Expanding Kustomize")
	}

	var list manifest.List
	d := yaml.NewDecoder(&buf)
	for {
		var m manifest.Manifest
		if err := d.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "Parsing Kustomize output")
		}

		// Kustomize might return "empty" elements in the YAML stream that consist
		// only of comments. Ignore these
		if len(m) == 0 {
			continue
		}

		list = append(list, m)
	}

	return list, nil
}
