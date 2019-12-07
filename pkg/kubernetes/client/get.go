package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Get retrieves a single Kubernetes object from the cluster
func (k Kubectl) Get(namespace, kind, name string) (manifest.Manifest, error) {
	return k.get(namespace, kind, []string{name})
}

// GetByLabels retrieves all objects matched by the given labels from the cluster
func (k Kubectl) GetByLabels(namespace, kind string, labels map[string]string) (manifest.List, error) {
	lArgs := make([]string, 0, len(labels))
	for k, v := range labels {
		lArgs = append(lArgs, fmt.Sprintf("-l=%s=%s", k, v))
	}

	list, err := k.get(namespace, kind, lArgs)
	if err != nil {
		return nil, err
	}

	if list.Kind() != "List" {
		return nil, fmt.Errorf("expected kind `List` but got `%s` instead", list.Kind())
	}

	items := list["items"].([]interface{})
	ms := make(manifest.List, 0, len(items))
	for _, i := range items {
		ms = append(ms, manifest.Manifest(i.(map[string]interface{})))
	}

	return ms, nil
}

func (k Kubectl) get(namespace, kind string, sel []string) (manifest.Manifest, error) {
	argv := []string{"get", "-o", "json"}
	if namespace == "" {
		argv = append(argv, "--all-namespaces")
	} else {
		argv = append(argv, "-n", namespace)
	}
	argv = append(argv, "--context", k.context.Get("name").MustStr())
	argv = append(argv, kind)
	argv = append(argv, sel...)

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

	var m manifest.Manifest
	if err := json.Unmarshal(sout.Bytes(), &m); err != nil {
		return nil, err
	}

	return m, nil
}
