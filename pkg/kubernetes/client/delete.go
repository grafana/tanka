package client

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Test-ability: isolate deleteCtl to build and return exec.Cmd from DeleteOpts
func (k Kubectl) deleteCtl(namespace, kind, name string, opts DeleteOpts) *exec.Cmd {
	argv := []string{
		"-n", namespace,
		kind, name,
	}
	if opts.Force {
		argv = append(argv, "--force")
	}

	if opts.DryRun != "" {
		dryRun := fmt.Sprintf("--dry-run=%s", opts.DryRun)
		argv = append(argv, dryRun)
	}

	return k.ctl("delete", argv...)
}

// Delete deletes the given Kubernetes resource from the cluster
func (k Kubectl) Delete(namespace, kind, name string, opts DeleteOpts) error {
	cmd := k.deleteCtl(namespace, kind, name, opts)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "Error from server (NotFound):") {
			print("Delete failed: " + stderr.String())
			return nil
		}
		return err
	}
	if opts.DryRun != "" {
		print(stdout.String())
	}

	return nil
}
