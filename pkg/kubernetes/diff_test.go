package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// TestSeparate checks that separate properly separates resources:
//
// - `namespace: ""` should be properly translated into the default namespace
// - cluster-wide resources are always `live`
// - resources with missing namespaces:
//   - `soon` if their namespace is included in the Jsonnet
//   - `live` otherwise, to cause a helpful error message
//     for the user that the namespace is indeed missing
func TestSeparate(t *testing.T) {
	cases := []struct {
		name string

		// state returned from Jsonnet
		state manifest.List
		// namespaces that exist
		namespaces []string

		// the default namespace (no value -> implicit default '')
		defaultNs string

		// resources that can be checked with the cluster
		live manifest.List
		// resources depending on a condition that will be met on next apply
		soon manifest.List
	}{
		{
			name:      "multi",
			defaultNs: "default",
			namespaces: []string{
				"default",
				"kube-system",
				"custom",
			},
			state: manifest.List{
				// cluster-wide resources: always live
				m("rbac.authorization.k8s.io/v1", "ClusterRole", "globalRole", ""),
				m("rbac.authorization.k8s.io/v1", "ClusterRoleBinding", "binding", "whydoihaveanamespace"),

				// default, existing namespace: `live`
				m("apps/v1", "Deployment", "loki", ""),
				m("apps/v1", "Deployment", "grafana", "default"),

				// custom, existing namespace: `live`
				m("apps/v1", "Deployment", "cortex", "custom"),

				// custom, soon existing namespace:
				m("v1", "Namespace", "monitoring", ""),                 // `live`
				m("apps/v1", "Deployment", "prometheus", "monitoring"), // `soon`

				// custom, missing namespace: `live`
				m("apps/v1", "Deployment", "metrictank", "metrics"),
			},
			live: manifest.List{
				m("rbac.authorization.k8s.io/v1", "ClusterRole", "globalRole", ""),
				m("rbac.authorization.k8s.io/v1", "ClusterRoleBinding", "binding", "whydoihaveanamespace"),
				m("apps/v1", "Deployment", "loki", ""),
				m("apps/v1", "Deployment", "grafana", "default"),
				m("apps/v1", "Deployment", "cortex", "custom"),
				m("v1", "Namespace", "monitoring", ""),
				m("apps/v1", "Deployment", "metrictank", "metrics"),
			},
			soon: manifest.List{
				m("apps/v1", "Deployment", "prometheus", "monitoring"),
			},
		},
		{
			name:       "default/soon",
			defaultNs:  "grafana",
			namespaces: []string{"default", "kube-system"}, // `grafana` missing
			state: manifest.List{
				m("v1", "Namespace", "grafana", ""), // `grafana` created during apply
				m("apps/v1", "Deployment", "prometheus", "grafana"),
				m("apps/v1", "Deployment", "cortex", ""), // implicit default `""`
			},
			live: manifest.List{
				m("v1", "Namespace", "grafana", ""),
			},
			soon: manifest.List{
				m("apps/v1", "Deployment", "prometheus", "grafana"),
				m("apps/v1", "Deployment", "cortex", ""),
			},
		},
		{
			name:       "default/missing",
			defaultNs:  "grafana",
			namespaces: []string{"default", "kube-system"}, // `grafana` missing
			state: manifest.List{
				// m("", "Namespace", "grafana", ""), <- `grafana` NOT created
				m("apps/v1", "Deployment", "prometheus", "grafana"),
				m("apps/v1", "Deployment", "cortex", ""), // implicit default (`""`)
			},
			live: manifest.List{
				// `live`, so user notices missing ns
				m("apps/v1", "Deployment", "prometheus", "grafana"),
				m("apps/v1", "Deployment", "cortex", ""),
			},
		},
	}

	// static set of resources for this test (usually obtained using
	// `client.Resources()`)
	staticResources := client.Resources{
		{APIVersion: "", Kind: "Namespace", Namespaced: false},
		{APIVersion: "apps/v1", Kind: "Deployment", Namespaced: true},
		{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole", Namespaced: false},
		{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRoleBinding", Namespaced: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			namespaces := make(map[string]bool)
			for _, n := range c.namespaces {
				namespaces[n] = true
			}

			live, soon := separate(c.state, c.defaultNs, separateOpts{
				namespaces: namespaces,
				resources:  staticResources,
			})

			assert.ElementsMatch(t, c.live, live, "live")
			assert.ElementsMatch(t, c.soon, soon, "soon")
		})
	}
}

func m(apiVersion, kind, name, namespace string) manifest.Manifest {
	return manifest.Manifest{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
	}
}
