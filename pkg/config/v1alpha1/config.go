package v1alpha1

// Config holds the configuration variables for config version v1alpha1
// ApiVersion and Kind are currently unused, this may change in the future.
type Config struct {
	APIVersion string       `json:"api_version"`
	Kind       string       `json:"kind"`
	Metadata   Metadata     `json:"metadata"`
	Spec       ProviderSpec `json:"spec"`
}

// Metadata is meant for humans and not parsed
type Metadata struct {
	Labels map[string]interface{} `json:"labels"`
}

// ProviderSpec is used to dynamically configure providers.
// The providers are iterated over and the first valid provider will be used.
// All properties specified inside of the provider are settings specific
// to the provider itself.
type ProviderSpec map[string]interface{}
