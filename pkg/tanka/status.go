package tanka

import (
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

type ParseResult struct {
	Env       *v1alpha1.Config
	Resources client.Manifests

	kube *kubernetes.Kubernetes
}

func Status(baseDir string, mods ...Modifier) (*ParseResult, error) {
	opts := parseModifiers(mods)

	sum, err := parse(baseDir, opts)
	if err != nil {
		return nil, err
	}

	return sum, nil
}
