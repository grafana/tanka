package client

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/objx"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Kubectl uses the `kubectl` command to operate on a Kubernetes cluster
type Kubectl struct {
	context   objx.Map
	cluster   objx.Map
	APIServer string
}

// New returns a instance of Kubectl with a correct context already discovered.
func New(endpoint string) (*Kubectl, error) {
	k := Kubectl{
		APIServer: endpoint,
	}
	if err := k.setupContext(); err != nil {
		return nil, errors.Wrap(err, "finding usable context")
	}
	return &k, nil
}

// Info returns known informational data about the client and its environment
func (k Kubectl) Info() (*Info, error) {
	client, server, err := k.version()
	if err != nil {
		return nil, errors.Wrap(err, "obtaining versions")
	}

	return &Info{
		ClientVersion: client,
		ServerVersion: server,

		Context: k.context,
		Cluster: k.cluster,
	}, nil
}

// Version returns the version of kubectl and the Kubernetes api server
func (k Kubectl) version() (client, server *semver.Version, err error) {
	cmd := exec.Command("kubectl", "version",
		"-o", "json",
		"--context", k.context.Get("name").MustStr(),
	)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}
	vs := objx.MustFromJSON(buf.String())
	client = semver.MustParse(vs.Get("clientVersion.gitVersion").MustStr())
	server = semver.MustParse(vs.Get("serverVersion.gitVersion").MustStr())
	return client, server, nil
}

// Diff takes a desired state as yaml and returns the differences
// to the system in common diff format
func (k Kubectl) DiffServerSide(data manifest.List) (*string, error) {
	argv := []string{"diff",
		"--context", k.context.Get("name").MustStr(),
		"-f", "-",
	}
	cmd := exec.Command("kubectl", argv...)

	raw := bytes.Buffer{}
	cmd.Stdout = &raw
	cmd.Stderr = FilterWriter{regexp.MustCompile(`exit status \d`)}

	cmd.Stdin = strings.NewReader(data.String())

	err := cmd.Run()

	// kubectl uses exit status 1 to tell us that there is a diff
	if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
		s := raw.String()
		return &s, nil
	}
	// another error
	if err != nil {
		return nil, err
	}

	// no diff -> nil
	return nil, nil
}

// FilterWriter is an io.Writer that discards every message that matches at
// least one of the regular expressions.
type FilterWriter []*regexp.Regexp

func (r FilterWriter) Write(p []byte) (n int, err error) {
	for _, re := range r {
		if re.Match(p) {
			// silently discard
			return len(p), nil
		}
	}
	return os.Stderr.Write(p)
}
