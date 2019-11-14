package tanka

import (
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

type Info struct {
	Env       *v1alpha1.Config
	Resources client.Manifests
	Client    client.Info
}

func Status(baseDir string, mods ...Modifier) (*Info, error) {
	opts := parseModifiers(mods)

	r, err := parse(baseDir, opts)
	if err != nil {
		return nil, err
	}
	r.Env.Spec.DiffStrategy = r.kube.Spec.DiffStrategy

	cInfo, err := r.kube.Info()
	if err != nil {
		return nil, err
	}

	return &Info{
		Env:       r.Env,
		Resources: r.Resources,
		Client:    *cInfo,
	}, nil
}
