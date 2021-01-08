package tanka

import (
	"log"

	"github.com/fatih/color"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
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
	DirLayout
}

type DirLayout struct {
	Base       string
	Root       string
	Entrypoint string
}

// Status returns information about the particular environment
func Status(path string, opts Opts) (*Info, error) {
	opts.EvalScript = MetaOnlyEvalScript

	root, base, err := jpath.Dirs(path)
	if err != nil {
		return nil, err
	}

	entry, err := jpath.Entrypoint(path)
	if err != nil {
		return nil, err
	}

	r, err := Load(path, opts)
	if err != nil {
		return nil, err
	}

	info := Info{
		Env:       r.Env,
		Resources: r.Resources,
		DirLayout: DirLayout{
			Root:       root,
			Base:       base,
			Entrypoint: entry,
		},
	}

	if r.Env.Spec.APIServer != "" {
		kube, err := r.Connect()
		if err != nil {
			log.Printf("%s %s\n\n", color.RedString("Error:"), err)
			return &info, nil
		}

		info.Env.Spec.DiffStrategy = kube.Env.Spec.DiffStrategy
		info.Client = kube.Info()
	}

	return &info, nil
}
