package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Kubectl uses the `kubectl` command to operate on a Kubernetes cluster
type Kubectl struct{}

// Get retrieves an Kubernetes object from the API
func (k Kubectl) Get(namespace, kind, name string) (map[string]interface{}, error) {
	argv := []string{"get", "-o", "json", "-n", namespace, kind, name}
	cmd := exec.Command("kubectl", argv...)
	raw := bytes.Buffer{}
	cmd.Stdout = &raw
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(raw.Bytes(), &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// Diff takes a desired state as yaml and returns the differences
// to the system in common diff format
func (k Kubectl) Diff(yaml string) (string, error) {
	argv := []string{"diff", "-f", "-"}
	cmd := exec.Command("kubectl", argv...)
	raw := bytes.Buffer{}
	cmd.Stdout = &raw

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		fmt.Fprintln(stdin, yaml)
		stdin.Close()
	}()

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// kubectl uses this to tell us that there is a diff
			if exitError.ExitCode() == 1 {
				return raw.String(), nil
			}
		}
		return "", err
	}

	return raw.String(), nil
}
