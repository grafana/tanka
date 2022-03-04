package v1alpha1

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// New creates a new Environment object with internal values already set
func New() *Environment {
	c := Environment{}

	// constants
	c.APIVersion = "tanka.dev/v1alpha1"
	c.Kind = "Environment"

	// default namespace
	c.Spec.Namespace = "default"

	c.Metadata.Labels = make(map[string]string)

	return &c
}

// Environment represents a set of resources in relation to its Kubernetes cluster
type Environment struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Metadata   Metadata    `json:"metadata"`
	Spec       Spec        `json:"spec"`
	Data       interface{} `json:"data,omitempty"`
}

// Metadata is meant for humans and not parsed
type Metadata struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// Has and Get make Metadata a simple wrapper for labels.Labels to use our map in their querier
func (m Metadata) Has(label string) (exists bool) {
	_, exists = m.Labels[label]
	return exists
}

// Get implements Get for labels.Labels interface
func (m Metadata) Get(label string) (value string) {
	return m.Labels[label]
}

func (m Metadata) NameLabel() string {
	partsHash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s", m.Name, m.Namespace)))
	chars := []rune(hex.EncodeToString(partsHash[:]))
	return string(chars[:48])
}

// Spec defines Kubernetes properties
type Spec struct {
	APIServer        string           `json:"apiServer"`
	ContextNames     []string         `json:"contextNames"`
	Namespace        string           `json:"namespace"`
	DiffStrategy     string           `json:"diffStrategy,omitempty"`
	InjectLabels     bool             `json:"injectLabels,omitempty"`
	ResourceDefaults ResourceDefaults `json:"resourceDefaults"`
	ExpectVersions   ExpectVersions   `json:"expectVersions"`
}

// ExpectVersions holds semantic version constraints
// TODO: extend this to handle more than Tanka
type ExpectVersions struct {
	Tanka string `json:"tanka,omitempty"`
}

// ResourceDefaults will be inserted in any manifests that tanka processes.
type ResourceDefaults struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}
