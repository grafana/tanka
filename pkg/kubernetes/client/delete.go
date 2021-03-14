package client

import (
	"bytes"
	"os"
	"strings"
)

func (k Kubectl) Delete(namespace, kind, name string, opts DeleteOpts) error {
	argv := []string{
		"-n", namespace,
		kind, name,
	}
	if opts.Force {
		argv = append(argv, "--force")
	}

	cmd := k.ctl("delete", argv...)
	// https://stackoverflow.com/questions/18159704/how-to-debug-exit-status-1-error-when-running-exec-command-in-golang

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

	return nil
}
