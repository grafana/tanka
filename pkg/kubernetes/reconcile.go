package kubernetes

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	funk "github.com/thoas/go-funk"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Reconcile extracts kubernetes Manifests from raw evaluated jsonnet <kind>/<name>,
// provided the manifests match the given regular expressions. It finds each manifest by
// recursively walking the jsonnet structure.
//
// In addition, we sort the manifests to ensure the order is consistent in each
// show/diff/apply cycle. This isn't necessary, but it does help users by producing
// consistent diffs.
func Reconcile(raw map[string]interface{}, spec v1alpha1.Spec, kindNameMatchers []*regexp.Regexp) (state manifest.List, err error) {
	manifests, err := walkJSON(raw)
	if err != nil {
		return nil, errors.Wrap(err, "flattening manifests")
	}

	// If we don't have a namespace, we want to set it to the default that is configured in
	// our kubernetes specification
	for _, manifest := range manifests {
		if spec.Namespace != "" && !manifest.Metadata().HasNamespace() {
			manifest.Metadata()["namespace"] = spec.Namespace
		}
	}

	// If we have any kind-name matchers, we should filter all the manifests by matching
	// against their <kind>/<name> identifier.
	if len(kindNameMatchers) > 0 {
		manifests = funk.Filter(manifests, func(elem interface{}) bool {
			manifest := elem.(manifest.Manifest)
			for _, matcher := range kindNameMatchers {
				if matcher.MatchString(manifest.KindName()) {
					return true
				}
			}

			return false
		}).([]manifest.Manifest)
	}

	sort.SliceStable(manifests, func(i int, j int) bool {
		return manifests[i].KindName() < manifests[j].KindName()
	})

	return manifests, nil
}

type ErrorPrimitiveReached struct {
	path, key string
	primitive interface{}
}

func (e ErrorPrimitiveReached) Error() string {
	return fmt.Sprintf("recursion did not resolve in a valid Kubernetes object. "+
		" In path path `%s` found key `%s` of type `%T` instead.",
		e.path, e.key, e.primitive)
}

// walkJSON recurses into either a map or list, returning a list of all objects that look
// like kubernetes resources. We support resources at an arbitrary level of nesting, and
// return an error if any leaf nodes f
//
// Handling the different types is quite gross, so we split this method into a generic
// walkJSON, and then walkObj/walkList to handle the two different types of collection we
// support.
func walkJSON(ptr interface{}, paths ...string) (manifest.List, error) {
	switch v := ptr.(type) {
	case map[string]interface{}:
		return walkObj(v, paths...)
	case []interface{}:
		return walkList(v, paths...)
	}

	key := paths[len(paths)-1]
	path := "." + strings.Join(paths[:len(paths)-1], ".")

	return nil, ErrorPrimitiveReached{path: path, key: key, primitive: ptr}
}

//		case d == nil: // result from false if condition in Jsonnet
func walkObj(obj objx.Map, paths ...string) (manifest.List, error) {
	obj = obj.Exclude([]string{"__ksonnet"}) // remove our private ksonnet field

	// This looks like a kubernetes manifest, so make one and return it
	if isKubernetesManifest(obj) {
		m := manifest.Manifest(obj.Value().MSI())

		return manifest.List{m}, nil
	}

	manifests := manifest.List{}
	for key, value := range obj {
		if value == nil { // result from false if condition in Jsonnet
			continue
		}
		children, err := walkJSON(value, append(paths, key)...)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, children...)
	}

	return manifests, nil
}

func walkList(list []interface{}, paths ...string) (manifest.List, error) {
	manifests := manifest.List{}
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
