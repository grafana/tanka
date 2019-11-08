package client

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Apply applies the given yaml to the cluster
func (k Kubectl) Apply(data Manifests, opts ApplyOpts) error {
	if err := k.setupContext(); err != nil {
		return errors.Wrap(err, "finding usable context")
	}

	argv := []string{"apply",
		"--context", k.context.Get("name").MustStr(),
		"-f", "-",
	}
	if opts.Force {
		argv = append(argv, "--force")
	}

	cmd := exec.Command("kubectl", argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Stdin = strings.NewReader(data.String())

	return cmd.Run()
}
