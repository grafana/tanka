package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/stretchr/objx"
	funk "github.com/thoas/go-funk"
)

// setupContext uses `kubectl config view` to obtain the KUBECONFIG and extracts the correct context from it
func (k *Kubectl) setupContext() error {
	if k.context != nil {
		return nil
	}

	var err error
	k.cluster, k.context, err = ContextFromIP(k.APIServer)
	if err != nil {
		return err
	}
	return nil
}

// Kubeconfig returns the merged $KUBECONFIG of the host
func Kubeconfig() (map[string]interface{}, error) {
	cmd := exec.Command("kubectl", "config", "view", "-o", "json")
	cfgJSON := bytes.Buffer{}
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(cfgJSON.Bytes(), &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Contexts() ([]string, error) {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o=name")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return strings.Split(buf.String(), "\n"), nil
}

// ContextFromIP searches the $KUBECONFIG for a context of a cluster that matches the apiServer
func ContextFromIP(apiServer string) (cluster, context objx.Map, err error) {
	kubeconfig, err := Kubeconfig()
	if err != nil {
		return nil, nil, err
	}
	cfg := objx.New(kubeconfig)

	// find the correct cluster
	cluster = find(cfg.Get("clusters").MustMSISlice(), "cluster.server", apiServer)
	if cluster == nil { // empty map means no result
		return nil, nil, fmt.Errorf("no cluster that matches the apiServer `%s` was found. Please check your $KUBECONFIG", apiServer)
	}

	// find a context that uses the cluster
	context = find(cfg.Get("contexts").MustMSISlice(), "context.cluster", cluster.Get("name").MustStr())
	if context == nil {
		return nil, nil, fmt.Errorf("no context that matches the cluster `%s` was found. Please check your $KUBECONFIG", cluster.Get("name").MustStr())
	}

	return cluster, context, nil
}

func IPFromContext(name string) (ip string, err error) {
	kubeconfig, err := Kubeconfig()
	if err != nil {
		return "", err
	}
	cfg := objx.New(kubeconfig)

	// find the context
	context := find(cfg.Get("contexts").MustMSISlice(), "name", name)
	if context == nil {
		return "", fmt.Errorf("no context named `%s` was found. Please check your $KUBECONFIG", name)
	}

	clusterName := context.Get("context.cluster").MustStr()
	cluster := find(cfg.Get("clusters").MustMSISlice(), "name", clusterName)
	if cluster == nil { // empty map means no result
		return "", fmt.Errorf("no cluster named `%s` as required by context `%s` was found. Please check your $KUBECONFIG", clusterName, name)
	}

	return cluster.Get("cluster.server").MustStr(), nil
}

func find(list []map[string]interface{}, prop string, expected string) objx.Map {
	return objx.New(funk.Find(list, func(x map[string]interface{}) bool {
		got := objx.New(x).Get(prop).MustStr()
		return got == expected
	}))
}
