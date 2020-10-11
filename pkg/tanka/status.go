package tanka

import (
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Info holds information about a particular environment, including its Config,
// the individual resources of the desired state and also the status of the
// client.
type Info struct {
	Env       *v1alpha1.Environment
	Resources manifest.List
	Client    client.Info
}

// Status returns information about the particular environment
func Status(path string, opts Opts) ([]*Info, error) {
	_, envs, err := ParseEnv(path, ParseOpts{JsonnetOpts: opts.JsonnetOpts})
	if err != nil {
		return nil, err
	}

	infos := make([]*Info, 0)
	for _, env := range envs {
		r, err := load(env, opts)
		if err != nil {
			return nil, err
		}
		kube, err := r.connect()
		if err != nil {
			return nil, err
		}

		r.Env.Spec.DiffStrategy = kube.Env.Spec.DiffStrategy

		info := &Info{
			Env:       r.Env,
			Resources: r.Resources,
			Client:    kube.Info(),
		}

		infos = append(infos, info)
	}
	return infos, nil
}
