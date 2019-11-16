package manifest

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	yaml "gopkg.in/yaml.v2"
)

// Manifest represents a Kubernetes API object. The fields `apiVersion` and
// `kind` are required, `metadata.name` should be present as well
type Manifest map[string]interface{}

func New(raw map[string]interface{}) (Manifest, error) {
	m := Manifest(raw)
	if err := m.Verify(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m Manifest) String() string {
	y, err := yaml.Marshal(m)
	if err != nil {
		// this should never go wrong in normal operations
		panic(errors.Wrap(err, "formatting manifest"))
	}
	return string(y)
}

// Verify checks whether the manifest is correctly structured
func (m Manifest) Verify() error {
	o := objx.New(map[string]interface{}(m))
	err := make(SchemaError)

	if !o.Has("kind") {
		err.add("kind")
	}
	if !o.Has("apiVersion") {
		err.add("apiVersion")
	}
	if !o.Has("metadata") {
		err.add("metadata")
	}
	if !o.Has("metadata.name") && m.Kind() != "List" {
		err.add("metadata.name")
	}

	if len(err) == 0 {
		return nil
	}
	return err
}

// Kind returns the kind of the API object
func (m Manifest) Kind() string {
	return m["kind"].(string)
}

// APIVersion returns the version of the API this object uses
func (m Manifest) APIVersion() string {
	return m["apiVersion"].(string)
}

// Metadata returns the metadata of this object
func (m Manifest) Metadata() Metadata {
	return Metadata(m["metadata"].(map[string]interface{}))
}

func (m Manifest) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	return m.Verify()
}

func (m Manifest) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&m); err != nil {
		return err
	}
	return m.Verify()
}

// Metadata is the metadata object from the Manifest
type Metadata map[string]interface{}

func (m Metadata) Name() string {
	return m["name"].(string)
}

func (m Metadata) HasNamespace() bool {
	return objx.New(m).Get("namespace").IsStr()
}
func (m Metadata) Namespace() string {
	if !m.HasNamespace() {
		return ""
	}
	return m["namespace"].(string)
}

func (m Metadata) HasLabels() bool {
	return objx.New(m).Get("labels").IsStr()
}
func (m Metadata) Labels() map[string]interface{} {
	if !m.HasLabels() {
		return make(map[string]interface{})
	}
	return m["labels"].(map[string]interface{})
}

func (m Metadata) HasAnnotations() bool {
	return objx.New(m).Get("annotations").IsStr()
}
func (m Metadata) Annotations() map[string]interface{} {
	if !m.HasAnnotations() {
		return make(map[string]interface{})
	}
	return m["annotations"].(map[string]interface{})
}

// Manifests is a list of individual Manifests
type List []Manifest

// String returns the Manifests as a yaml stream. In case of an error, it is
// returned as a string instead.
func (m List) String() string {
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)

	for _, d := range m {
		if err := enc.Encode(d); err != nil {
			// This should never happen in normal operations
			panic(errors.Wrap(err, "formatting manifests"))
		}
	}

	return buf.String()
}
