package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/stretchr/objx"
	funk "github.com/thoas/go-funk"
)

// setupContext uses `kubectl config view` to obtain the KUBECONFIG and extracts the correct context from it
func (k *Kubectl) setupContext() error {
	if k.context != nil {
		return nil
	}

	cmd := exec.Command("kubectl", "config", "view", "-o", "json")
	cfgJSON := bytes.Buffer{}
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(cfgJSON.Bytes(), &cfg); err != nil {
		return err
	}

	var err error
	k.cluster, k.context, err = contextFromKubeconfig(cfg, k.APIServer)
	if err != nil {
		return err
	}
	return nil
}

// contextFromKubeconfig searches a kubeconfig for a context of a cluster that matches the apiServer
func contextFromKubeconfig(kubeconfig map[string]interface{}, apiServer string) (cluster, context objx.Map, err error) {
	cfg := objx.New(kubeconfig)

	// find the correct cluster
	cluster = objx.New(funk.Find(cfg.Get("clusters").MustMSISlice(), func(x map[string]interface{}) bool {
		host := objx.New(x).Get("cluster.server").MustStr()
		return host == apiServer
	}))
	if !(len(cluster) > 0) { // empty map means no result
		return nil, nil, fmt.Errorf("no cluster that matches the apiServer `%s` was found. Please check your $KUBECONFIG", apiServer)
	}

	// find a context that uses the cluster
	context = objx.New(funk.Find(cfg.Get("contexts").MustMSISlice(), func(x map[string]interface{}) bool {
		c := objx.New(x)
		return c.Get("context.cluster").MustStr() == cluster.Get("name").MustStr()
	}))
	if !(len(context) > 0) {
		return nil, nil, fmt.Errorf("no context that matches the cluster `%s` was found. Please check your $KUBECONFIG", cluster.Get("name").MustStr())
	}

	return cluster, context, nil
}
