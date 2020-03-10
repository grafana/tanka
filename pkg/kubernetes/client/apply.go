package client

import (
	"os"
	"strings"

	funk "github.com/thoas/go-funk"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Apply applies the given yaml to the cluster
func (k Kubectl) Apply(data manifest.List, opts ApplyOpts) error {
	// create namespaces first to succeed first try
	ns := filterNamespace(data)
	if len(ns) > 0 {
		if err := k.apply(ns, opts); err != nil {
			return err
		}
	}

	return k.apply(data, opts)
}

func (k Kubectl) apply(data manifest.List, opts ApplyOpts) error {
	argv := []string{"-f", "-"}
	if opts.Force {
		argv = append(argv, "--force")
	}

	if !opts.Validate {
		argv = append(argv, "--validate=false")
	}

	cmd := k.ctl("apply", argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Stdin = strings.NewReader(data.String())

	return cmd.Run()
}

func filterNamespace(in manifest.List) manifest.List {
	return manifest.List(funk.Filter(in, func(i manifest.Manifest) bool {
		return strings.ToLower(i.Kind()) == "namespace"
	}).([]manifest.Manifest))
}
