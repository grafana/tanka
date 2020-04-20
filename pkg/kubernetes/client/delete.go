package client

import (
	"os"
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
