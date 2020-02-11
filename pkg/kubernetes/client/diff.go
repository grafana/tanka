package client

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/util"
)

// DiffServerSide takes the desired state and computes the differences on the
// server, returning them in `diff(1)` format
func (k Kubectl) DiffServerSide(data manifest.List) (*string, error) {
	ns, err := k.Namespaces()
	if err != nil {
		return nil, err
	}

	ready, missing := separateMissingNamespace(data, ns)
	cmd := k.ctl("diff", "-f", "-")

	raw := bytes.Buffer{}
	cmd.Stdout = &raw
	cmd.Stderr = FilterWriter{regexp.MustCompile(`exit status \d`)}

	cmd.Stdin = strings.NewReader(ready.String())

	err = cmd.Run()

	// kubectl uses exit status 1 to tell us that there is a diff
	if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
	} else if err != nil {
		return nil, err
	}

	s := raw.String()
	for _, r := range missing {
		d, err := util.DiffStr(util.DiffName(r), "", r.String())
		if err != nil {
			return nil, err
		}
		s += d
	}

	if s != "" {
		return &s, nil
	}

	// no diff -> nil
	return nil, nil
}

func separateMissingNamespace(in manifest.List, exists map[string]bool) (ready, missingNamespace manifest.List) {
	for _, r := range in {
		if !exists[r.Metadata().Namespace()] {
			missingNamespace = append(missingNamespace, r)
			continue
		}
		ready = append(ready, r)
	}
	return
}
