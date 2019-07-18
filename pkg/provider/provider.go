package provider

import "github.com/spf13/cobra"

// Provider describes methods for functionality more advanced than evaluating jsonnet
type Provider interface {
	// Show receives the raw evaluated jsonnet as a marshaled json dict and
	// shall return it reconciled as a state object of the target system
	Show(raw map[string]interface{}) (state map[string]interface{}, err error)

	// Apply receives a state object generated using `Show()` and may apply it to the target system
	Apply(desired map[string]interface{}) error

	// State shall return the current state of the target system.
	// It receives the desired state object generated using `Show()`.
	// This is used for diffing afterwards.
	State(desired map[string]interface{}) (real map[string]interface{}, err error)

	// Cmd shall return a command to be available under `tk provider <NAME>`
	Cmd() *cobra.Command
}

// EmptyConstructor defines a function interface that creates uninitialized Providers,
// ready to be unmarshalled from Config
type EmptyConstructor func() Provider
