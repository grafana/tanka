package client

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// APIResources retrieves a list of supported API resources from a Kubernetes API server
func (k Kubectl) APIResources(opts APIResourcesOpts) ([]string, error) {
	argv := append([]string{"api-resources"})
	if len(opts.Verbs) > 0 {
		argv = append(argv, "--verbs", strings.Join(opts.Verbs, ","))
	}
	if opts.Output != "" {
		argv = append(argv, "-o", opts.Output)
	}

	cmd := exec.Command("kubectl", argv...)

	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr

	if err := cmd.Run(); err != nil {
		if strings.HasPrefix(serr.String(), "Error from server (NotFound)") {
			return nil, ErrorNotFound{serr.String()}
		}
		if strings.HasPrefix(serr.String(), "error: the server doesn't have a resource type") {
			return nil, ErrorUnknownResource{serr.String()}
		}

		fmt.Print(serr.String())
		return nil, err
	}

	s := strings.TrimSpace(sout.String())
	return strings.Split(s, "\n"), nil
}
