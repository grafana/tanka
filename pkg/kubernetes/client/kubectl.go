package client

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// Namespaces of the cluster
func (k Kubectl) Namespaces() (map[string]bool, error) {
	argv := []string{"get",
		"-o", "json",
		"--context", k.context.Get("name").MustStr(),
		"namespaces",
	}
	cmd := exec.Command("kubectl", argv...)

	var sout bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var list manifest.Manifest
	if err := json.Unmarshal(sout.Bytes(), &list); err != nil {
		return nil, err
	}

	items, ok := list["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("listing namespaces: expected items to be an object, but got %T instead", list["items"])
	}

	namespaces := make(map[string]bool)
	for _, i := range items {
		m, err := manifest.New(i.(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		namespaces[m.Metadata().Name()] = true
	}
	return namespaces, nil
}

func (k Kubectl) APIResources() ([]string, error) {
	argv := []string{"api-resources",
		"--context", k.context.Get("name").MustStr(),
		"--verbs=list", "-o=name",
	}
	cmd := exec.Command("kubectl", argv...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n"), nil
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
