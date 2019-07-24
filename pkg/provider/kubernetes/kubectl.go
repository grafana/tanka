package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/stretchr/objx"
	funk "github.com/thoas/go-funk"
)

// Kubectl uses the `kubectl` command to operate on a Kubernetes cluster
type Kubectl struct {
	context   string
	APIServer string
}

// setupContext uses `kubectl config view` to obtain the KUBECONFIG and extracts the correct context from it
func (k Kubectl) setupContext() error {
	cmd := exec.Command("kubectl", "config", "view", "-o", "json")
	cfgJSON := bytes.Buffer{}
	cmd.Stdout = &cfgJSON
	if err := cmd.Run(); err != nil {
		return err
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(cfgJSON.Bytes(), &cfg); err != nil {
		return err
	}

	var err error
	k.context, err = contextFromKubeconfig(cfg, k.APIServer)
	if err != nil {
		return err
	}
	return nil
}

// contextFromKubeconfig searches a kubeconfig for a context of a cluster that matches the apiServer
func contextFromKubeconfig(kubeconfig map[string]interface{}, apiServer string) (string, error) {
	cfg := objx.New(kubeconfig)

	// find the correct cluster
	cluster := objx.New(funk.Find(cfg.Get("clusters").MustMSISlice(), func(x map[string]interface{}) bool {
		host := objx.New(x).Get("cluster.server").MustStr()
		return host == apiServer
	}))
	if !(len(cluster) > 0) { // empty map means no result
		return "", fmt.Errorf("no cluster that matches the apiServer `%s` was found. Please check your $KUBECONFIG", apiServer)
	}

	// find a context that uses the cluster
	context := objx.New(funk.Find(cfg.Get("contexts").MustMSISlice(), func(x map[string]interface{}) bool {
		c := objx.New(x)
		return c.Get("context.cluster").MustStr() == cluster.Get("name").MustStr()
	}))
	if !(len(context) > 0) {
		return "", fmt.Errorf("no context that matches the cluster `%s` was found. Please check your $KUBECONFIG", cluster.Get("name").MustStr())
	}

	return context.Get("name").MustStr(), nil
}

// Get retrieves an Kubernetes object from the API
func (k Kubectl) Get(namespace, kind, name string) (map[string]interface{}, error) {
	if err := k.setupContext(); err != nil {
		return nil, err
	}
	argv := []string{"get",
		"-o", "json",
		"-n", namespace,
		"--context", k.context,
		kind, name,
	}
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
	if err := k.setupContext(); err != nil {
		return "", err
	}
	argv := []string{"diff",
		"--context", k.context,
		"-f", "-",
	}
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
