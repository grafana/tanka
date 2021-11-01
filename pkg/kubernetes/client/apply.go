package client

import (
	"os"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Apply applies the given yaml to the cluster
func (k Kubectl) Apply(data manifest.List, opts ApplyOpts) error {
	argv := []string{"-f", "-"}
	if opts.Force {
		argv = append(argv, "--force")
	}

	if !opts.Validate {
		argv = append(argv, "--validate=false")
	}

	if !opts.InsecureSkipTlsVerify {
		argv = append(argv, "--insecure-skip-tls-verify")
	}

	cmd := k.ctl("apply", argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Stdin = strings.NewReader(data.String())

	return cmd.Run()
}
