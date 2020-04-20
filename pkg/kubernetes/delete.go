package kubernetes

import (
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

type DeleteOpts client.DeleteOpts

func (k *Kubernetes) Delete(state manifest.List, opts DeleteOpts) error {
	for _, m := range state {
		if err := k.ctl.Delete(m.Metadata().Namespace(), m.Kind(), m.Metadata().Name(), client.DeleteOpts(opts)); err != nil {
			return err
		}
	}

	return nil
}
