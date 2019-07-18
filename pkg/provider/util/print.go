package util

import (
	"bytes"
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

// ShowJSON marshals the state into plain JSON
func ShowJSON(state interface{}) (string, error) {
	out, err := json.MarshalIndent(state, "", "  ")
	return string(out), err
}

// ShowYAML marshals the state into a single YAML document.
func ShowYAML(state interface{}) (string, error) {
	out, err := yaml.Marshal(state)
	return string(out), err
}

// ShowYAMLDocs the state into multiple yaml documents
func ShowYAMLDocs(state []map[string]interface{}) (string, error) {
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)

	for _, d := range state {
		if err := enc.Encode(d); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
