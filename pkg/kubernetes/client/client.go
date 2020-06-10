package client

import (
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/resources"
)

// Client for working with Kubernetes
type Client interface {
	// Get the specified object(s) from the cluster
	Get(namespace, kind, name string) (manifest.Manifest, error)
	GetByLabels(namespace, kind string, labels map[string]string) (manifest.List, error)
	GetByState(data manifest.List) (manifest.List, error)

	// Apply the configuration to the cluster. `data` must contain a plaintext
	// format that is `kubectl-apply(1)` compatible
	Apply(data manifest.List, opts ApplyOpts) error

	// DiffServerSide runs the diff operation on the server and returns the
	// result in `diff(1)` format
	DiffServerSide(data manifest.List) (*string, error)

	// Delete the specified object(s) from the cluster
	Delete(namespace, kind, name string, opts DeleteOpts) error

	// Namespaces the cluster currently has
	Namespaces() (map[string]bool, error)
	// Resources returns all known api-resources of the cluster
	Resources() (resources.Store, error)

	// Info returns known informational data about the client. Best effort based,
	// fields of `Info` that cannot be stocked with valuable data, e.g.
	// due to an error, shall be left nil.
	Info() Info

	// Close may run tasks once the client is no longer needed.
	Close() error
}

// ApplyOpts allow to specify additional parameter for apply operations
type ApplyOpts struct {
	// force allows to ignore checks and force the operation
	Force bool

	// validate allows to enable/disable kubectl validation
	Validate bool

	// autoApprove allows to skip the interactive approval
	AutoApprove bool
}

// DeleteOpts allow to specify additional parameters for delete operations
// Currently not different from ApplyOpts, but may be required in the future
type DeleteOpts ApplyOpts
