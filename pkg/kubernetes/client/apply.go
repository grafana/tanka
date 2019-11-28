package client

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	funk "github.com/thoas/go-funk"
)

// Apply applies the given yaml to the cluster
func (k Kubectl) Apply(labels []string, data manifest.List, opts ApplyOpts) error {
	// create namespaces first to succeed first try
	ns := filterNamespace(data)
	if err := k.apply(labels, ns, opts); err != nil {
		return err
	}

	return k.apply(labels, data, opts)
}

func (k Kubectl) apply(labels []string, data manifest.List, opts ApplyOpts) error {
	argv := []string{"apply",
		"--context", k.context.Get("name").MustStr(),
		"-f", "-",
	}
	if opts.Force {
		argv = append(argv, "--force")
	}

	if opts.Prune {
		labelString := strings.Join(labels, ",")
		argv = append(argv, "--prune", "-l", labelString)
	}

	cmd := exec.Command("kubectl", argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Stdin = strings.NewReader(data.String())

	return cmd.Run()
}

// ApplyDryRun reports the changes that are due to be applied, including prunings
func (k Kubectl) ApplyDryRun(labels []string, data manifest.List) (string, error) {

	labelString := strings.Join(labels, ",")
	argv := []string{"apply",
		"--context", k.context.Get("name").MustStr(),
		"-f", "-",
		"--prune", "-l", labelString,
		"--dry-run",
	}

	cmd := exec.Command("kubectl", argv...)
	raw := bytes.Buffer{}
	cmd.Stdout = &raw
	cmd.Stderr = os.Stderr

	cmd.Stdin = strings.NewReader(data.String())

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return raw.String(), nil
}

func filterNamespace(in manifest.List) manifest.List {
	return manifest.List(funk.Filter(in, func(i manifest.Manifest) bool {
		return strings.ToLower(i.Kind()) == "namespace"
	}).([]manifest.Manifest))
}
