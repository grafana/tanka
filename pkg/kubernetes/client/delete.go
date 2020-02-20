package client

import (
	"fmt"
	"os"
)

// Delete removes the specified object from the cluster
func (k Kubectl) Delete(namespace, kind, name string, opts DeleteOpts) error {
	return k.delete(namespace, []string{kind, name}, opts)
}

// DeleteByLabels removes all objects matched by the given labels from the cluster
func (k Kubectl) DeleteByLabels(namespace string, labels map[string]interface{}, opts DeleteOpts) error {
	lArgs := make([]string, 0, len(labels))
	for k, v := range labels {
		lArgs = append(lArgs, fmt.Sprintf("-l=%s=%s", k, v))
	}

	return k.delete(namespace, lArgs, opts)
}

func (k Kubectl) delete(namespace string, sel []string, opts DeleteOpts) error {
	argv := append([]string{"-n", namespace}, sel...)
	k.ctl("delete", argv...)

	if opts.Force {
		argv = append(argv, "--force")
	}

	cmd := KubectlCmd(argv...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
