package client

import (
	"testing"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/stretchr/testify/assert"
)

func TestSeparateMissingNamespace(t *testing.T) {
	cases := []struct {
		name string
		td   nsTd

		missing bool
	}{
		// default should always exist
		{
			name: "default",
			td: newNsTd(func(m manifest.Metadata) {
				m["namespace"] = "default"
			}, []string{}),
			missing: false,
		},
		// implcit default (not specfiying an ns at all) also
		{
			name: "implicit-default",
			td: newNsTd(func(m manifest.Metadata) {
				delete(m, "namespace")
			}, []string{}),
			missing: false,
		},
		// custom ns that exists
		{
			name: "custom-ns",
			td: newNsTd(func(m manifest.Metadata) {
				m["namespace"] = "custom"
			}, []string{"custom"}),
			missing: false,
		},
		// custom ns that does not exist
		{
			name: "missing-ns",
			td: newNsTd(func(m manifest.Metadata) {
				m["namespace"] = "missing"
			}, []string{}),
			missing: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ready, missing := separateMissingNamespace(manifest.List{c.td.m}, c.td.ns)
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

func manifestOfKind(kind, apiVersion string) manifest.Manifest {
	return manifest.Manifest{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata":   map[string]interface{}{},
	}
}

func TestSeparateUnknownResources(t *testing.T) {
	cases := []struct {
		name    string
		in      manifest.List
		kinds   map[string]bool
		ok      manifest.List
		missing manifest.List
	}{
		{
			name:    "empty",
			in:      manifest.List{},
			kinds:   map[string]bool{"Pod": true},
			ok:      manifest.List{},
			missing: manifest.List{},
		},
		{
			name:    "all-supported",
			in:      manifest.List{manifestOfKind("Pod", "v1")},
			kinds:   map[string]bool{"Pod": true},
			ok:      manifest.List{manifestOfKind("Pod", "v1")},
			missing: manifest.List{},
		},
		{
			name:    "one-unsupported",
			in:      manifest.List{manifestOfKind("Pod", "v1"), manifestOfKind("Custom", "abcd/v1")},
			kinds:   map[string]bool{"Pod": true},
			ok:      manifest.List{manifestOfKind("Pod", "v1")},
			missing: manifest.List{manifestOfKind("Custom", "abcd/v1")},
		},
		{
			name:    "all-unsupported",
			in:      manifest.List{manifestOfKind("Pod", "v1"), manifestOfKind("Custom", "abcd/v1")},
			kinds:   map[string]bool{"Deployment": true},
			ok:      manifest.List{},
			missing: manifest.List{manifestOfKind("Pod", "v1"), manifestOfKind("Custom", "abcd/v1")},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ready, missing := separateUnknownResources(c.in, c.kinds)
			assert.ElementsMatch(t, ready, c.ok)
			assert.ElementsMatch(t, missing, c.missing)
		})
	}
}
