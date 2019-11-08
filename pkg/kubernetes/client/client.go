package client

import (
	"github.com/Masterminds/semver"
	"github.com/stretchr/objx"
)

type Client interface {
	Get(namespace, kind, name string) (Manifest, error)
	GetByLabels(namespace string, labels map[string]interface{}) (Manifests, error)

	// Apply the configuration to the cluster. `data` must contain a plaintext
	// format that is `kubectl-apply(1)` compatible
	Apply(data Manifests, opts ApplyOpts) error

	// DiffServerSide runs the diff operation on the server and returns the
	// result in `diff(1)` format
	DiffServerSide(data Manifests) (*string, error)

	// Delete the specified object from the cluster
	Delete(namespace, kind, name string, opts DeleteOpts) error
	DeleteByLabels(namespace string, labels map[string]interface{}, opts DeleteOpts) error

	// Info returns known informational data about the client. Best effort based,
	// fields of `Info` that cannot be stocked with valuable data, e.g.
	// due to an error, shall be left nil.
	Info() (*Info, error)
}

// Info contains metadata about the client and its environment
type Info struct {
	// version of `kubectl`
	ClientVersion *semver.Version

	// version of the API server
	ServerVersion *semver.Version

	// used context and cluster from KUBECONFIG
	Context, Cluster objx.Map
}

// ApplyOpts allow to specify additional parameter for apply operations
type ApplyOpts struct {
	// force allows to ignore checks and force the operation
	Force bool

	// autoApprove allows to skip the interactive approval
	AutoApprove bool
}

// DeleteOpts allow to specify additional parameters for delete operations
// Currently not different from ApplyOpts, but may be required in the future
type DeleteOpts ApplyOpts
