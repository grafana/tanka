package provider

import "github.com/spf13/cobra"

// Provider describes methods for functionality more advanced than evaluating jsonnet
type Provider interface {
	// Init may be used to do some setup related tasks.
	// Called before any other requests to the provider are made
	Init() error

	// Reconcile receives the raw evaluated jsonnet as a marshaled json dict and
	// shall return it reconciled as a state object of the target system
	// `state` must be one of {`map[string]interface`, `[]map[string]interface`}
	Reconcile(raw map[string]interface{}) (state interface{}, err error)

	// Diff receives the state object generated using `Reconcile()`
	// and may return the differences to the current state.
	Diff(state interface{}) (string, error)

	// Apply receives a state object generated using `Reconcile()`
	// and may apply it to the target system
	Apply(desired interface{}) error

	// Fmt receives the state object generated using `Reconcile()`
	// and may pretty-print it into the string.
	Fmt(state interface{}) (string, error)

	// Cmd shall return a command to be available under `tk provider <NAME>`
	Cmd() *cobra.Command
}

// EmptyConstructor defines a function interface that creates uninitialized Providers,
// ready to be unmarshalled from Config
type EmptyConstructor func() Provider
