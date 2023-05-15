package helm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/rs/zerolog/log"
)

// Helm provides high level access to some Helm operations
type Helm interface {
	// Pull downloads a Helm Chart from a remote
	Pull(chart, version string, opts PullOpts) error

	// RepoUpdate fetches the latest remote index
	RepoUpdate(opts Opts) error

	// Template returns the individual resources of a Helm Chart
	Template(name, chart string, opts TemplateOpts) (manifest.List, error)

	// ChartExists checks if a chart exists in the provided calledFromPath
	ChartExists(chart string, opts *JsonnetOpts) (string, error)
}

// PullOpts are additional, non-required options for Helm.Pull
type PullOpts struct {
	Opts

	// Directory to put the resulting .tgz into
	Destination string

	// Where to extract the chart to, defaults to the name of the chart
	ExtractDirectory string
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

	// Pull to a temp dir within the destination directory (not /tmp) to avoid possible cross-device issues when renaming
	tempDir, err := os.MkdirTemp(opts.Destination, ".pull-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	chartPullPath := chart
	chartRepoName, chartName := parseReqRepo(chart), parseReqName(chart)
	for _, configuredRepo := range opts.Repositories {
		if configuredRepo.Name == chartRepoName {
			// OCI images are pulled with their full path
			if strings.HasPrefix(configuredRepo.URL, "oci://") {
				chartPullPath = fmt.Sprintf("%s/%s", configuredRepo.URL, chartName)
			}
		}
	}

	cmd := e.cmd("pull", chartPullPath,
		"--version", version,
		"--repository-config", repoFile,
		"--destination", tempDir,
		"--untar",
	)

	if err = cmd.Run(); err != nil {
		return err
	}

	if opts.ExtractDirectory == "" {
		opts.ExtractDirectory = chartName
	}

	// It is not possible to tell `helm pull` to extract to a specific directory
	// so we extract to a temp dir and then move the files to the destination
	return os.Rename(
		filepath.Join(tempDir, chartName),
		filepath.Join(opts.Destination, opts.ExtractDirectory),
	)
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

func (e ExecHelm) ChartExists(chart string, opts *JsonnetOpts) (string, error) {
	// resolve the Chart relative to the caller
	callerDir := filepath.Dir(opts.CalledFrom)
	chart = filepath.Join(callerDir, chart)
	if _, err := os.Stat(chart); err != nil {
		return "", fmt.Errorf("helmTemplate: Failed to find a chart at '%s': %s. See https://tanka.dev/helm#failed-to-find-chart", chart, err)
	}

	return chart, nil
}

// cmd returns a prepared exec.Cmd to use the `helm` binary
func (e ExecHelm) cmd(action string, args ...string) *exec.Cmd {
	argv := []string{action}
	argv = append(argv, args...)
	log.Debug().Strs("argv", argv).Msg("running helm")

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

	f, err := os.CreateTemp("", "charts-repos")
	if err != nil {
		return "", err
	}

	enc := json.NewEncoder(f)
	if err := enc.Encode(m); err != nil {
		return "", err
	}

	return f.Name(), nil
}
