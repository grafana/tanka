package provider

import "github.com/spf13/cobra"

// Provider describes methods for functionality more advanced than evaluating jsonnet
type Provider interface {
	// Reconcile receives the raw evaluated jsonnet as a marshaled json dict and
	// shall return it reconciled as a state object of the target system
	// `state` must be one of {`map[string]interface`, `[]map[string]interface`}
	Reconcile(raw map[string]interface{}) (state interface{}, err error)


	// Apply receives a state object generated using `Reconcile()`
	// and may apply it to the target system
	Apply(desired interface{}) error

	// State shall return the current state of the target system.
	// It receives the desired state object generated using `Format()`.
	// This is used for diffing afterwards.
	State(desired interface{}) (real map[string]interface{}, err error)
	// Fmt receives the state object generated using `Reconcile()`
	// and may pretty-print it into the string.
	Fmt(state interface{}) (string, error)

	// Cmd shall return a command to be available under `tk provider <NAME>`
	Cmd() *cobra.Command
}

// EmptyConstructor defines a function interface that creates uninitialized Providers,
// ready to be unmarshalled from Config
type EmptyConstructor func() Provider
