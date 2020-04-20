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

// Order in which install different kinds of Kubernetes objects.
// Inspired by https://github.com/helm/helm/blob/8c84a0bc0376650bc3d7334eef0c46356c22fa36/pkg/releaseutil/kind_sorter.go
var kindOrder = []string{
	"Namespace",
	"NetworkPolicy",
	"ResourceQuota",
	"LimitRange",
	"PodSecurityPolicy",
	"PodDisruptionBudget",
	"ServiceAccount",
	"Secret",
	"ConfigMap",
	"StorageClass",
	"PersistentVolume",
	"PersistentVolumeClaim",
	"CustomResourceDefinition",
	"ClusterRole",
	"ClusterRoleList",
	"ClusterRoleBinding",
	"ClusterRoleBindingList",
	"Role",
	"RoleList",
	"RoleBinding",
	"RoleBindingList",
	"Service",
	"DaemonSet",
	"Pod",
	"ReplicationController",
	"ReplicaSet",
	"Deployment",
	"HorizontalPodAutoscaler",
	"StatefulSet",
	"Job",
	"CronJob",
	"Ingress",
	"APIService",
}

const (
	MetadataPrefix   = "tanka.dev"
	LabelEnvironment = MetadataPrefix + "/environment"
)

// Reconcile extracts kubernetes Manifests from raw evaluated jsonnet <kind>/<name>,
// provided the manifests match the given regular expressions. It finds each manifest by
// recursively walking the jsonnet structure.
//
// In addition, we sort the manifests to ensure the order is consistent in each
// show/diff/apply cycle. This isn't necessary, but it does help users by producing
// consistent diffs.
func Reconcile(raw map[string]interface{}, cfg v1alpha1.Config, targets []*regexp.Regexp) (state manifest.List, err error) {
	extracted, err := extract(raw)
	if err != nil {
		return nil, errors.Wrap(err, "flattening manifests")
	}

	out := make(manifest.List, 0, len(extracted))
	for _, m := range extracted {

		// inject tanka.dev/environment label
		if cfg.Spec.InjectLabels {
			m.Metadata().Labels()[LabelEnvironment] = cfg.Metadata.NameLabel()
		}

		out = append(out, m)
	}

	// If we have any kind-name matchers, we should filter all the manifests by matching
	// against their <kind>/<name> identifier.
	if len(targets) > 0 {
		tmp := funk.Filter(out, func(i interface{}) bool {
			p := objectspec(i.(manifest.Manifest))
			for _, t := range targets {
				if t.MatchString(p) {
					return true
				}
			}
			return false
		}).([]manifest.Manifest)
		out = manifest.List(tmp)
	}

	// Stable output order
	sort.SliceStable(out, func(i int, j int) bool {
		var io, jo int

		// anything that is not in kindOrder will get to the end of the install list.
		for io = 0; io < len(kindOrder); io++ {
			if out[i].Kind() == kindOrder[io] {
				break
			}
		}

		for jo = 0; jo < len(kindOrder); jo++ {
			if out[j].Kind() == kindOrder[jo] {
				break
			}
		}

		// If Kind of both objects are at different indexes of kindOrder, sort by them
		if io != jo {
			return io < jo
		}

		// If the Kinds themselves are different (e.g. both of the Kinds are not in
		// the kindOrder), order alphabetically.
		if out[i].Kind() != out[j].Kind() {
			return out[i].Kind() < out[j].Kind()
		}

		// If namespaces differ, sort by the namespace.
		if out[i].Metadata().Namespace() != out[j].Metadata().Namespace() {
			return out[i].Metadata().Namespace() < out[j].Metadata().Namespace()
		}

		// Otherwise, order the objects by name.
		return out[i].Metadata().Name() < out[j].Metadata().Name()
	})

	return out, nil
}

func extract(deep interface{}) (map[string]manifest.Manifest, error) {
	extracted := make(map[string]manifest.Manifest)
	if err := walkJSON(deep, extracted, nil); err != nil {
		return nil, err
	}
	return extracted, nil
}

// walkJSON recurses into either a map or list, returning a list of all objects that look
// like kubernetes resources. We support resources at an arbitrary level of nesting, and
// return an error if a node is not walkable.
//
// Handling the different types is quite gross, so we split this method into a generic
// walkJSON, and then walkObj/walkList to handle the two different types of collection we
// support.
func walkJSON(ptr interface{}, extracted map[string]manifest.Manifest, path trace) error {
	// check for known types
	switch v := ptr.(type) {
	case map[string]interface{}:
		return walkObj(v, extracted, path)
	case []interface{}:
		return walkList(v, extracted, path)
	}

	return ErrorPrimitiveReached{
		path:      path.Base(),
		key:       path.Name(),
		primitive: ptr,
	}
}

func walkList(list []interface{}, extracted map[string]manifest.Manifest, path trace) error {
	for idx, value := range list {
		err := walkJSON(value, extracted, append(path, fmt.Sprintf("[%d]", idx)))
		if err != nil {
			return err
		}
	}
	return nil
}

func walkObj(obj objx.Map, extracted map[string]manifest.Manifest, path trace) error {
	obj = obj.Exclude([]string{"__ksonnet"}) // remove our private ksonnet field

	// This looks like a kubernetes manifest, so make one and return it
	if isKubernetesManifest(obj) {
		m, err := manifest.NewFromObj(obj)
		if err != nil {
			return err.(*manifest.SchemaError).WithName(path.Full())
		}

		extracted[path.Full()] = m
		return nil
	}

	for key, value := range obj {
		path := append(path, key)

		if value == nil { // result from false if condition in Jsonnet
			continue
		}
		err := walkJSON(value, extracted, path)
		if err != nil {
			return err
		}
	}

	return nil
}

type trace []string

func (t trace) Full() string {
	return "." + strings.Join(t, ".")
}

func (t trace) Base() string {
	if len(t) > 0 {
		t = t[:len(t)-1]
	}
	return "." + strings.Join(t, ".")
}

func (t trace) Name() string {
	if len(t) > 0 {
		return t[len(t)-1]
	}

	return ""
}

// ErrorPrimitiveReached occurs when walkJSON reaches the end of nested dicts without finding a valid Kubernetes manifest
type ErrorPrimitiveReached struct {
	path, key string
	primitive interface{}
}

func (e ErrorPrimitiveReached) Error() string {
	return fmt.Sprintf("recursion did not resolve in a valid Kubernetes object. "+
		" In path `%s` found key `%s` of type `%T` instead.",
		e.path, e.key, e.primitive)
}

// isKubernetesManifest attempts to infer whether the given object is a valid kubernetes
// resource by verifying the presence of apiVersion and kind. These two
// fields are required for kubernetes to accept any resource.
//
// In future, it might be a good idea to allow users to opt their object out of being
// interpreted as a kubernetes resource, perhaps with a field like `exclude: true`. For
// now, any object within the jsonnet output that quacks like a kubernetes resource will
// be provided to the kubernetes API.
func isKubernetesManifest(obj objx.Map) bool {
	return obj.Get("apiVersion").IsStr() && obj.Get("apiVersion").Str() != "" &&
		obj.Get("kind").IsStr() && obj.Get("kind").Str() != ""
}
