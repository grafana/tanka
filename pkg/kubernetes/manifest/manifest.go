package manifest

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	yaml "gopkg.in/yaml.v2"
)

// Manifest represents a Kubernetes API object. The fields `apiVersion` and
// `kind` are required, `metadata.name` should be present as well
type Manifest map[string]interface{}

// New creates a new Manifest
func New(raw map[string]interface{}) (Manifest, error) {
	m := Manifest(raw)
	if err := m.Verify(); err != nil {
		return nil, err
	}
	return m, nil
}

// NewFromObj creates a new Manifest from an objx.Map
func NewFromObj(raw objx.Map) (Manifest, error) {
	return New(map[string]interface{}(raw))
}

// String returns the Manifest in yaml representation
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
	o := m2o(m)
	var err SchemaError

	if !o.Get("kind").IsStr() {
		err.add("kind")
	}
	if !o.Get("apiVersion").IsStr() {
		err.add("apiVersion")
	}

	// Lists don't have `metadata`
	if !strings.HasSuffix(m.Kind(), "List") {
		if !o.Get("metadata").IsMSI() {
			err.add("metadata")
		}
		if !o.Get("metadata.name").IsStr() {
			err.add("metadata.name")
		}
	}

	if len(err.fields) == 0 {
		return nil
	}

	return &err
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
	if m["metadata"] == nil {
		m["metadata"] = make(map[string]interface{})
	}
	return Metadata(m["metadata"].(map[string]interface{}))
}

// UnmarshalJSON validates the Manifest during json parsing
func (m *Manifest) UnmarshalJSON(data []byte) error {
	type tmp Manifest
	var t tmp
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*m = Manifest(t)
	return m.Verify()
}

// UnmarshalYAML validates the Manifest during yaml parsing
func (m *Manifest) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tmp Manifest
	var t tmp
	if err := unmarshal(&t); err != nil {
		return err
	}
	*m = Manifest(t)
	return m.Verify()
}

// Metadata is the metadata object from the Manifest
type Metadata map[string]interface{}

// Name of the manifest
func (m Metadata) Name() string {
	name, ok := m["name"]
	if !ok {
		return ""
	}
	return name.(string)
}

// HasNamespace returns whether the manifest has a namespace set
func (m Metadata) HasNamespace() bool {
	return m2o(m).Get("namespace").IsStr()
}

// Namespace of the manifest
func (m Metadata) Namespace() string {
	if !m.HasNamespace() {
		return ""
	}
	return m["namespace"].(string)
}

// HasLabels returns whether the manifest has labels
func (m Metadata) HasLabels() bool {
	return m2o(m).Get("labels").IsMSI()
}

// Labels of the manifest
func (m Metadata) Labels() map[string]interface{} {
	if !m.HasLabels() {
		return make(map[string]interface{})
	}
	return m["labels"].(map[string]interface{})
}

func (m Metadata) SetLabel(key, value string) {
	if !m.HasLabels() {
		m["labels"] = make(map[string]interface{})
	}
	m["labels"].(map[string]interface{})[key] = value
}

// HasAnnotations returns whether the manifest has annotations
func (m Metadata) HasAnnotations() bool {
	return m2o(m).Get("annotations").IsMSI()
}

// Annotations of the manifest
func (m Metadata) Annotations() map[string]interface{} {
	if !m.HasAnnotations() {
		return make(map[string]interface{})
	}
	return m["annotations"].(map[string]interface{})
}

// List of individual Manifests
type List []Manifest

// String returns the List as a yaml stream. In case of an error, it is
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

func m2o(m interface{}) objx.Map {
	switch mm := m.(type) {
	case Metadata:
		return objx.New(map[string]interface{}(mm))
	case Manifest:
		return objx.New(map[string]interface{}(mm))
	}
	return nil
}
