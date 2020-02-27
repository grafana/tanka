package client

import (
	"testing"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/stretchr/testify/assert"
)

func TestSeparateMissingNamespace(t *testing.T) {
	cases := []struct {
		name      string
		td        nsTd
		namespace string // namespace as defined in spec.json

		missing bool
	}{
		// default should always exist
		{
			name:      "default",
			namespace: "default",
			td: newNsTd(func(m manifest.Metadata) {
				m["namespace"] = "default"
			}, []string{}),
			missing: false,
		},
		// implicit default (not specfiying an ns at all) also
		{
			name:      "implicit-default",
			namespace: "default",
			td: newNsTd(func(m manifest.Metadata) {
				delete(m, "namespace")
			}, []string{}),
			missing: false,
		},
		// custom ns that exists
		{
			name:      "custom-ns",
			namespace: "custom",
			td: newNsTd(func(m manifest.Metadata) {
				m["namespace"] = "custom"
			}, []string{"custom"}),
			missing: false,
		},
		// custom ns that does not exist
		{
			name:      "missing-ns",
			namespace: "missing",
			td: newNsTd(func(m manifest.Metadata) {
				m["namespace"] = "missing"
			}, []string{}),
			missing: true,
		},
		// an explicitly created namespace is missing
		{
			name:      "explicit-namespace-missing",
			namespace: "other-namespace",
			td: newNsTd(func(m manifest.Metadata) {
				m["kind"] = "Namespace"
				m["name"] = "explicit-namespace"
			}, []string{}),
			missing: true,
		},
		// an environment's namespace is missing and resources have no ns defined
		{
			name:      "env-namespace-missing",
			namespace: "env-namespace",
			td:        newNsTd(func(m manifest.Metadata) {}, []string{}),
			missing:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ready, missing := separateMissingNamespace(manifest.List{c.td.m}, c.namespace, c.td.ns)
			if c.missing {
				assert.Lenf(t, ready, 0, "expected manifest to be missing (ready = 0)")
				assert.Lenf(t, missing, 1, "expected manifest to be missing (missing = 1)")
			} else {
				assert.Lenf(t, ready, 1, "expected manifest to be ready (ready = 1)")
				assert.Lenf(t, missing, 0, "expected manifest to be ready (missing = 0)")
			}
		})
	}
}

type nsTd struct {
	m  manifest.Manifest
	ns map[string]bool
}

func newNsTd(f func(m manifest.Metadata), ns []string) nsTd {
	m := manifest.Manifest{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata":   map[string]interface{}{},
	}
	if f != nil {
		f(m.Metadata())
	}

	nsMap := map[string]bool{
		"default": true, // you can't get rid of this one ever
	}
	for _, n := range ns {
		nsMap[n] = true
	}

	return nsTd{
		m:  m,
		ns: nsMap,
	}
}
