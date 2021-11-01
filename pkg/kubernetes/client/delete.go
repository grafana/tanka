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

        if !opts.InsecureSkipTlsVerify {
                argv = append(argv, "--insecure-skip-tls-verify")
        }

	cmd := k.ctl("delete", argv...)

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
