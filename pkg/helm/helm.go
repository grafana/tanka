package helm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Helm provides high level access to some Helm operations
type Helm interface {
	// Pull downloads a Helm Chart from a remote
	Pull(chart, version string, opts PullOpts) error

	// RepoUpdate fetches the latest remote index
	RepoUpdate(opts Opts) error

	// Template returns the individual resources of a Helm Chart
	Template(name, chart string, opts TemplateOpts) (manifest.List, error)
}

// PullOpts are additional, non-required options for Helm.Pull
type PullOpts struct {
	Opts

	// Directory to put the resulting .tgz into
	Destination string
}

// Opts are additional, non-required options that all Helm operations accept
type Opts struct {
	Repositories []Repo
}

// ExecHelm is a Helm implementation powered by the `helm` command line utility
type ExecHelm struct{}

// Pull implements Helm.Pull
func (e ExecHelm) Pull(chart, version string, opts PullOpts) error {
	repoFile, err := writeRepoTmpFile(opts.Repositories)
	if err != nil {
		return err
	}
	defer os.Remove(repoFile)

	cmd := e.cmd("pull", chart,
		"--version", version,
		"--destination", opts.Destination,
		"--repository-config", repoFile,
		"--untar",
	)

	return cmd.Run()
}

// RepoUpdate implements Helm.RepoUpdate
func (e ExecHelm) RepoUpdate(opts Opts) error {
	repoFile, err := writeRepoTmpFile(opts.Repositories)
	if err != nil {
		return err
	}
	defer os.Remove(repoFile)

	cmd := e.cmd("repo", "update",
		"--repository-config", repoFile,
	)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s\n%s", errBuf.String(), err)
	}

	return nil
}

// cmd returns a prepared exec.Cmd to use the `helm` binary
func (e ExecHelm) cmd(action string, args ...string) *exec.Cmd {
	argv := []string{action}
	argv = append(argv, args...)

	cmd := helmCmd(argv...)
	cmd.Stderr = os.Stderr

	return cmd
}

// helmCmd returns a bare exec.Cmd pointed at the local helm binary
func helmCmd(args ...string) *exec.Cmd {
	bin := "helm"
	if env := os.Getenv("TANKA_HELM_PATH"); env != "" {
		bin = env
	}

	return exec.Command(bin, args...)
}

// writeRepoTmpFile creates a temporary repositories.yaml from the passed Repo
// slice to be used by the helm binary
func writeRepoTmpFile(r []Repo) (string, error) {
	m := map[string]interface{}{
		"repositories": r,
	}

	f, err := ioutil.TempFile("", "charts-repos")
	if err != nil {
		return "", err
	}

	enc := json.NewEncoder(f)
	if err := enc.Encode(m); err != nil {
		return "", err
	}

	return f.Name(), nil
}
