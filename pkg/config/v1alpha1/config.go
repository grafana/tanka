package v1alpha1

type Config struct {
	APIVersion string       `json:"api_version"`
	Kind       string       `json:"kind"`
	Metadata   Metadata     `json:"metadata"`
	Spec       ProviderSpec `json:"spec"`
}

type Metadata struct {
	Labels map[string]interface{} `json:"labels"`
}

type ProviderSpec map[string]interface{}
