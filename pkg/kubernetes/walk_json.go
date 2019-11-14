package kubernetes

import (
	"fmt"
	"strings"

	"github.com/stretchr/objx"
)

// Manifest describes a single Kubernetes manifest. A manifest should only be constructed
// if the source data has been validated with isKubernetesManifest.
type Manifest map[string]interface{}

func (m Manifest) APIVersion() string {
	return m["apiVersion"].(string)
}

func (m Manifest) Kind() string {
	return m["kind"].(string)
}

func (m Manifest) Name() string {
	return m["metadata"].(map[string]interface{})["name"].(string)
}

func (m Manifest) Namespace() string {
	return m["metadata"].(map[string]interface{})["namespace"].(string)
}

func (m Manifest) KindName() string {
	return fmt.Sprintf("%s/%s", m.Kind(), m.Name())
}

func (m Manifest) GroupVersionKindName() string {
	return fmt.Sprintf("%s/%s/%s", m.APIVersion(), m.Kind(), m.Name())
}

func (m Manifest) Get(key string) *objx.Value {
	return m.Objx().Get(key)
}

func (m Manifest) Set(key string, value interface{}) Manifest {
	return Manifest(m.Objx().Set(key, value))
}

func (m Manifest) Objx() objx.Map {
	return objx.New((map[string]interface{})(m))
}

// walkJSON recurses into either a map or list, returning a list of all objects that look
// like kubernetes resources. We support resources at an arbitrary level of nesting, and
// return an error if any leaf nodes f
//
// Handling the different types is quite gross, so we split this method into a generic
// walkJSON, and then walkObj/walkList to handle the two different types of collection we
// support.
func walkJSON(ptr interface{}, paths ...string) ([]Manifest, error) {
	if obj, ok := ptr.(map[string]interface{}); ok {
		return walkObj(obj, paths...)
	}

	if list, ok := ptr.([]interface{}); ok {
		return walkList(list, paths...)
	}

	return nil, fmt.Errorf(
		"found object at %s that does not look like a Kubernetes object, or list",
		strings.Join(paths, "/"),
	)
}

func walkObj(obj objx.Map, paths ...string) ([]Manifest, error) {
	obj = obj.Exclude([]string{"__ksonnet"}) // remove our private ksonnet field

	// This looks like a kubernetes manifest, so make one and return it
	if isKubernetesManifest(obj) {
		manifest := Manifest(obj.Value().MSI())

		return []Manifest{manifest}, nil
	}

	manifests := []Manifest{}
	for key, value := range obj {
		children, err := walkJSON(value, append(paths, key)...)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, children...)
	}

	return manifests, nil
}

func walkList(list []interface{}, paths ...string) ([]Manifest, error) {
	manifests := []Manifest{}
	for idx, value := range list {
		children, err := walkJSON(value, append(paths, fmt.Sprintf("%d", idx))...)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, children...)
	}

	return manifests, nil
}

// isKubernetesManifest attempts to infer whether the given object is a valid kubernetes
// resource by verifying the presence of apiVersion, kind and metadata.name. These three
// fields are required for kubernetes to accept any resource.
//
// In future, it might be a good idea to allow users to opt their object out of being
// interpreted as a kubernetes resource, perhaps with a field like `exclude: true`. For
// now, any object within the jsonnet output that quacks like a kubernetes resource will
// be provided to the kubernetes API.
func isKubernetesManifest(obj objx.Map) bool {
	return obj.Get("apiVersion").IsStr() && obj.Get("apiVersion").Str() != "" &&
		obj.Get("kind").IsStr() && obj.Get("kind").Str() != "" &&
		obj.Get("metadata.name").IsStr() && obj.Get("metadata.name").Str() != ""
}
