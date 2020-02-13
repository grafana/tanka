package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	"github.com/stretchr/objx"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Kubectl uses the `kubectl` command to operate on a Kubernetes cluster
type Kubectl struct {
	// kubeconfig
	nsPatch string
	context objx.Map
	cluster objx.Map

	info *Info

	APIServer string
}

// New returns a instance of Kubectl with a correct context already discovered.
func New(endpoint, namespace string) (*Kubectl, error) {
	k := Kubectl{
		APIServer: endpoint,
	}
	if err := k.setupContext(namespace); err != nil {
		return nil, errors.Wrap(err, "finding usable context")
	}

	info, err := k.newInfo()
	if err != nil {
		return nil, errors.Wrap(err, "gathering client info")
	}
	k.info = info

	return &k, nil
}

// Info returns known informational data about the client and its environment
func (k Kubectl) Info() Info {
	return *k.info
}

func (k Kubectl) newInfo() (*Info, error) {
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
	cmd := k.ctl("version", "-o", "json")
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
	cmd := k.ctl("get", "namespaces", "-o", "json")

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

// FilterWriter is an io.Writer that discards every message that matches at
// least one of the regular expressions.
type FilterWriter struct {
	buf     string
	filters []*regexp.Regexp
}

func (r *FilterWriter) Write(p []byte) (n int, err error) {
	for _, re := range r.filters {
		if re.Match(p) {
			// silently discard
			return len(p), nil
		}
	}
	r.buf += string(p)
	return os.Stderr.Write(p)
}
