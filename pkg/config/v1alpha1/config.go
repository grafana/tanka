package v1alpha1

// Config holds the configuration variables for config version v1alpha1
// ApiVersion and Kind are currently unused, this may change in the future.
type Config struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

// Metadata is meant for humans and not parsed
type Metadata struct {
	Labels map[string]interface{} `json:"labels"`
}

// Spec defines Kubernetes properties
type Spec struct {
	APIServer    string `json:"apiServer"`
	Namespace    string `json:"namespace"`
	DiffStrategy string `json:"diffStrategy"`
}
