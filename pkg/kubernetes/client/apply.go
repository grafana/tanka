package client

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Test-ability: isolate applyCtl to build and return exec.Cmd from ApplyOpts
func (k Kubectl) applyCtl(data manifest.List, opts ApplyOpts) *exec.Cmd {
	argv := []string{"-f", "-"}
	if opts.Force {
		argv = append(argv, "--force")
	}

	if !opts.Validate {
		argv = append(argv, "--validate=false")
	}

	if opts.DryRun != "" {
		dryRun := fmt.Sprintf("--dry-run=%s", opts.DryRun)
		argv = append(argv, dryRun)
	}

	if opts.ServerSide {
		argv = append(argv, "--server-side")
	}

	return k.ctl("apply", argv...)
}

// Apply applies the given yaml to the cluster
func (k Kubectl) Apply(data manifest.List, opts ApplyOpts) error {
	cmd := k.applyCtl(data, opts)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Stdin = strings.NewReader(data.String())

	return cmd.Run()
}
