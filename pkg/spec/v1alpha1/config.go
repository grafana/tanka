package v1alpha1

import "strings"

// New creates a new Config object with internal values already set
func New() *Config {
	c := Config{}

	// constants
	c.APIVersion = "tanka.dev/v1alpha1"
	c.Kind = "Environment"

	// defaults
	c.Spec.Namespace = "default"
	c.Spec.InjectLabels.Environment = true

	c.Metadata.Labels = make(map[string]string)

	return &c
}

// Config holds the configuration variables for config version v1alpha1
// ApiVersion and Kind are currently unused, this may change in the future.
type Config struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

// Metadata is meant for humans and not parsed Name is an exception to this
// rule, as it's value is discovered at runtime from the directories name.
type Metadata struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

func (m Metadata) NameLabel() string {
	return strings.Replace(m.Name, "/", ".", -1)
}

// Spec defines Kubernetes properties
type Spec struct {
	APIServer    string       `json:"apiServer"`
	Namespace    string       `json:"namespace"`
	DiffStrategy string       `json:"diffStrategy,omitempty"`
	InjectLabels InjectLabels `json:"injectLabels,omitempty"`
}

// InjectLabels defines labels that Tanka will inject into manifests
type InjectLabels struct {
	Environment bool `json:"environment,omitempty"`
}
