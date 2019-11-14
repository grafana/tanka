package client

import (
	"bytes"
	"errors"

	yaml "gopkg.in/yaml.v2"
)

// Manifest represents a Kubernetes API object. The fields `apiVersion` and
// `kind` are required, `metadata.name` should be present as well
type Manifest map[string]interface{}

// Kind returns the kind of the API object
func (m Manifest) Kind() string {
	if _, ok := m["kind"]; !ok {
		return ""
	}
	return m["kind"].(string)
}

// APIVersion returns the version of the API this object uses
func (m Manifest) APIVersion() string {
	if _, ok := m["apiVersion"]; !ok {
		return ""
	}
	return m["apiVersion"].(string)
}

// Metadata returns the metadata of this object
func (m Manifest) Metadata() Metadata {
	if _, ok := m["metadata"]; !ok {
		return nil
	}
	return Metadata(m["metadata"].(map[string]interface{}))
}

// VerifyLax checks whether `kind`, `apiVersion` and `metadata` are present.
// Unless for some special type (e.g. List), Verify should be used instead.
func (m Manifest) VerifyLax() error {
	if m.Kind() == "" {
		return errors.New("kind missing")
	}
	if m.APIVersion() == "" {
		return errors.New("apiVersion missing")
	}
	if m.Metadata() == nil {
		return errors.New("metadata missing")
	}
	return nil
}

// Verify checks the same as VerifyLax but also requires `metadata.name` to be
// present.
func (m Manifest) Verify() error {
	if err := m.VerifyLax(); err != nil {
		return err
	}

	if m.Metadata().Name() == "" {
		return errors.New("name missing")
	}

	if m.Metadata().Labels() == nil {
		return errors.New("labels missing")
	}
	return nil
}

type Metadata map[string]interface{}

func (m Metadata) Name() string {
	if _, ok := m["name"]; !ok {
		return ""
	}
	return m["name"].(string)
}

func (m Metadata) Namespace() string {
	if _, ok := m["namespace"]; !ok {
		return ""
	}
	return m["namespace"].(string)
}

func (m Metadata) Labels() map[string]interface{} {
	if _, ok := m["labels"]; !ok {
		return nil
	}
	return m["labels"].(map[string]interface{})
}

type Manifests []Manifest

// String returns the Manifests as a yaml stream. In case of an error, it is
// returned as a string instead.
func (m Manifests) String() string {
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)

	for _, d := range m {
		if err := enc.Encode(d); err != nil {
			return err.Error()
		}
	}

	return buf.String()
}
