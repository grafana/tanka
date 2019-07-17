package kubernetes

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Kubernetes struct {
	apiServer string
	namespace string
}

// Show receives the raw evaluated jsonnet as a marshaled json dict and
// shall return it reconciled as a state object of the target system
func (k *Kubernetes) Show(raw map[string]interface{}) (state map[string]interface{}, err error) {
	panic("not implemented")
}

// Apply receives a state object generated using `Show()` and may apply it to the target system
func (k *Kubernetes) Apply(desired map[string]interface{}) error {
	panic("not implemented")
}

// State shall return the current state of the target system.
// It receives the desired state object generated using `Show()`.
// This is used for diffing afterwards.
func (k *Kubernetes) State(desired map[string]interface{}) (real map[string]interface{}, err error) {
	panic("not implemented")
}

// Cmd shall return a command to be available under `tk provider`
func (k *Kubernetes) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kubernetes",
		Short: "Kubernetes provider commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("not implemented")
		},
	}
}
