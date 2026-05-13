package kubernetes

import (
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// fakeClient is a minimal client.Client for testing Orphaned.
type fakeClient struct {
	resources    client.Resources
	byLabels     manifest.List
	byState      manifest.List
	byStateErr   error
	byLabelsKind string // records the kind string passed to GetByLabels
}

func (f *fakeClient) Get(namespace, kind, name string) (manifest.Manifest, error) {
	return nil, nil
}
func (f *fakeClient) Resources() (client.Resources, error) { return f.resources, nil }

func (f *fakeClient) GetByLabels(namespace, kind string, labels map[string]string) (manifest.List, error) {
	f.byLabelsKind = kind
	return f.byLabels, nil
}

func (f *fakeClient) GetByState(data manifest.List, opts client.GetByStateOpts) (manifest.List, error) {
	return f.byState, f.byStateErr
}

func (f *fakeClient) Apply(data manifest.List, opts client.ApplyOpts) error { return nil }
func (f *fakeClient) DiffServerSide(data manifest.List) (*string, error)    { return nil, nil }
func (f *fakeClient) DiffExitCode(data manifest.List) (bool, error)         { return false, nil }
func (f *fakeClient) Delete(namespace, apiVersion, kind, name string, opts client.DeleteOpts) error {
	return nil
}
func (f *fakeClient) Namespaces() (map[string]bool, error)         { return nil, nil }
func (f *fakeClient) Namespace(ns string) (manifest.Manifest, error) { return nil, nil }
func (f *fakeClient) Info() client.Info {
	return client.Info{ClientVersion: semver.MustParse("1.22.0")}
}
func (f *fakeClient) Close() error { return nil }

// testEnv returns a minimal environment suitable for Orphaned tests.
func testEnv() v1alpha1.Environment {
	return v1alpha1.Environment{
		Metadata: v1alpha1.Metadata{Name: "default/test"},
		Spec: v1alpha1.Spec{
			InjectLabels: true,
		},
	}
}

// orphanedManifest builds a manifest that looks as though it was directly
// applied by Tanka: it carries the last-applied annotation.
func orphanedManifest(kind, name, uid string) manifest.Manifest {
	return manifest.Manifest{
		"apiVersion": "apps/v1",
		"kind":       kind,
		"metadata": map[string]any{
			"name":      name,
			"namespace": "default",
			"uid":       uid,
			"annotations": map[string]any{
				AnnotationLastApplied: "{}",
			},
		},
	}
}

// appsResources returns a standard set of apps/v1 API resources for tests.
// Name is the plural resource name used to build the FQN (e.g. "statefulsets.apps").
func appsResources() client.Resources {
	return client.Resources{
		{APIVersion: "apps/v1", Kind: "StatefulSet", Name: "statefulsets", Namespaced: true, Verbs: "list,get"},
		{APIVersion: "apps/v1", Kind: "Deployment", Name: "deployments", Namespaced: true, Verbs: "list,get"},
	}
}

func TestOrphaned_NoFilters(t *testing.T) {
	sts := orphanedManifest("StatefulSet", "old-store", "uid-1")
	dep := orphanedManifest("Deployment", "old-app", "uid-2")

	fc := &fakeClient{
		resources: appsResources(),
		byState:   manifest.List{},
		byLabels:  manifest.List{sts, dep},
	}

	k := &Kubernetes{Env: testEnv(), ctl: fc}
	orphaned, err := k.Orphaned(manifest.List{}, OrphanedOpts{})
	require.NoError(t, err)

	assert.Len(t, orphaned, 2)
	// With no filters both FQNs are included in the cluster query.
	assert.Contains(t, fc.byLabelsKind, "statefulsets.apps")
	assert.Contains(t, fc.byLabelsKind, "deployments.apps")
}

func TestOrphaned_FilterByKind(t *testing.T) {
	sts := orphanedManifest("StatefulSet", "old-store", "uid-1")
	dep := orphanedManifest("Deployment", "old-app", "uid-2")

	fc := &fakeClient{
		resources: appsResources(),
		byState:   manifest.List{},
		// The fake always returns byLabels regardless of the kind argument; the
		// final process.Filter call in Orphaned ensures only StatefulSets are
		// returned to the caller.
		byLabels: manifest.List{sts, dep},
	}

	filters, err := process.StrExps("statefulset/.*")
	require.NoError(t, err)

	k := &Kubernetes{Env: testEnv(), ctl: fc}
	orphaned, err := k.Orphaned(manifest.List{}, OrphanedOpts{Filters: filters})
	require.NoError(t, err)

	// Only the StatefulSet passes the final filter.
	require.Len(t, orphaned, 1)
	assert.Equal(t, "StatefulSet", orphaned[0].Kind())

	// Kind pre-filtering: only the StatefulSet FQN was passed to GetByLabels.
	assert.Contains(t, fc.byLabelsKind, "statefulsets.apps")
	assert.NotContains(t, fc.byLabelsKind, "deployments.apps")
}

func TestOrphaned_FilterByKindAndName(t *testing.T) {
	live := orphanedManifest("StatefulSet", "live-store-0", "uid-1")
	old := orphanedManifest("StatefulSet", "old-store", "uid-2")

	fc := &fakeClient{
		resources: appsResources(),
		byState:   manifest.List{},
		byLabels:  manifest.List{live, old},
	}

	filters, err := process.StrExps("statefulset/old-.*")
	require.NoError(t, err)

	k := &Kubernetes{Env: testEnv(), ctl: fc}
	orphaned, err := k.Orphaned(manifest.List{}, OrphanedOpts{Filters: filters})
	require.NoError(t, err)

	// Only old-store matches the name pattern.
	require.Len(t, orphaned, 1)
	assert.Equal(t, "old-store", orphaned[0].Metadata().Name())
}

func TestOrphaned_KnownResourcesNotReturned(t *testing.T) {
	// known is still present in desired state; orphaned is not.
	known := orphanedManifest("StatefulSet", "live-store", "uid-keep")
	orphan := orphanedManifest("StatefulSet", "dead-store", "uid-drop")

	fc := &fakeClient{
		resources: client.Resources{
			{APIVersion: "apps/v1", Kind: "StatefulSet", Name: "statefulsets", Namespaced: true, Verbs: "list,get"},
		},
		byState:  manifest.List{known},
		byLabels: manifest.List{known, orphan},
	}

	k := &Kubernetes{Env: testEnv(), ctl: fc}
	orphaned, err := k.Orphaned(manifest.List{known}, OrphanedOpts{})
	require.NoError(t, err)

	require.Len(t, orphaned, 1)
	assert.Equal(t, "dead-store", orphaned[0].Metadata().Name())
}

func TestOrphaned_NegativeFilter(t *testing.T) {
	sts := orphanedManifest("StatefulSet", "old-store", "uid-1")
	dep := orphanedManifest("Deployment", "old-app", "uid-2")

	fc := &fakeClient{
		resources: appsResources(),
		byState:   manifest.List{},
		byLabels:  manifest.List{sts, dep},
	}

	// Negative filter: exclude StatefulSets from prune candidates.
	filters, err := process.StrExps("!statefulset/.*")
	require.NoError(t, err)

	k := &Kubernetes{Env: testEnv(), ctl: fc}
	orphaned, err := k.Orphaned(manifest.List{}, OrphanedOpts{Filters: filters})
	require.NoError(t, err)

	// Deployment survives; StatefulSet is excluded by the negative filter.
	require.Len(t, orphaned, 1)
	assert.Equal(t, "Deployment", orphaned[0].Kind())
	// No positive matcher → all kinds were queried.
	assert.Contains(t, fc.byLabelsKind, "statefulsets.apps")
	assert.Contains(t, fc.byLabelsKind, "deployments.apps")
}
