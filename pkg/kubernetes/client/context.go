package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/stretchr/objx"
	funk "github.com/thoas/go-funk"
)

// findContext returns a valid context from $KUBECONFIG that uses the given
// apiServer endpoint.
func findContext(endpoint string) (Config, error) {
	cluster, context, err := ContextFromIP(endpoint)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Context: *context,
		Cluster: *cluster,
	}, nil
}

// writeNamespacePatch writes a temporary file that includes only the previously
// discovered context with the `context.namespace` field set to the default
// namespace from `spec.json`. Adding this file to `$KUBECONFIG` results in
// `kubectl` picking this up, effectively setting the default namespace.
func writeNamespacePatch(context Context, defaultNamespace string) (string, error) {
	context.Context.Namespace = defaultNamespace

	patch := map[string]interface{}{
		"contexts": []interface{}{context},
	}
	out, err := json.Marshal(patch)
	if err != nil {
		return "", err
	}

	f, err := ioutil.TempFile("", "tk-kubectx-namespace-*.yaml")
	if err != nil {
		return "", err
	}
	if err = ioutil.WriteFile(f.Name(), out, 0644); err != nil {
		return "", err
	}

	return f.Name(), nil
}

// Kubeconfig returns the merged $KUBECONFIG of the host
func Kubeconfig() (objx.Map, error) {
	cmd := kubectlCmd("config", "view", "-o", "json")
	cfgJSON := bytes.Buffer{}
	cmd.Stdout = &cfgJSON
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return objx.FromJSON(cfgJSON.String())
}

// Contexts returns a list of context names
func Contexts() ([]string, error) {
	cmd := kubectlCmd("config", "get-contexts", "-o=name")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return strings.Split(buf.String(), "\n"), nil
}

// ContextFromIP searches the $KUBECONFIG for a context using a cluster that matches the apiServer
func ContextFromIP(apiServer string) (*Cluster, *Context, error) {
	cfg, err := Kubeconfig()
	if err != nil {
		return nil, nil, err
	}

	// find the correct cluster
	var cluster Cluster
	clusters, err := tryMSISlice(cfg.Get("clusters"), "clusters")
	if err != nil {
		return nil, nil, err
	}

	err = find(clusters, "cluster.server", apiServer, &cluster)
	if err == ErrorNoMatch {
		return nil, nil, ErrorNoCluster(apiServer)
	} else if err != nil {
		return nil, nil, err
	}

	// find a context that uses the cluster
	var context Context
	contexts, err := tryMSISlice(cfg.Get("contexts"), "contexts")
	if err != nil {
		return nil, nil, err
	}

	err = find(contexts, "context.cluster", cluster.Name, &context)
	if err == ErrorNoMatch {
		return nil, nil, ErrorNoContext(cluster.Name)
	} else if err != nil {
		return nil, nil, err
	}

	return &cluster, &context, nil
}

// IPFromContext parses $KUBECONFIG, finds the cluster with the given name and
// returns the cluster's endpoint
func IPFromContext(name string) (ip string, err error) {
	cfg, err := Kubeconfig()
	if err != nil {
		return "", err
	}

	// find a context with the given name
	var context Context
	contexts, err := tryMSISlice(cfg.Get("contexts"), "contexts")
	if err != nil {
		return "", err
	}

	err = find(contexts, "name", name, &context)
	if err == ErrorNoMatch {
		return "", ErrorNoContext(name)
	} else if err != nil {
		return "", err
	}

	// find the cluster of the context
	var cluster Cluster
	clusters, err := tryMSISlice(cfg.Get("clusters"), "clusters")
	if err != nil {
		return "", err
	}

	clusterName := context.Context.Cluster
	err = find(clusters, "name", clusterName, &cluster)
	if err == ErrorNoMatch {
		return "", fmt.Errorf("no cluster named `%s` as required by context `%s` was found. Please check your $KUBECONFIG", clusterName, name)
	} else if err != nil {
		return "", err
	}

	return cluster.Cluster.Server, nil
}

func tryMSISlice(v *objx.Value, what string) ([]map[string]interface{}, error) {
	if s := v.MSISlice(); s != nil {
		return s, nil
	}

	data, ok := v.Data().([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected %s to be of type `[]map[string]interface{}`, but got `%T` instead", what, v.Data())
	}
	return data, nil
}

// ErrorNoMatch occurs when no item matched had the expected value
var ErrorNoMatch error = errors.New("no matches found")

// find attempts to find an object in list whose prop equals expected.
// If found, the value is unmarshalled to ptr, otherwise errNotFound is returned.
func find(list []map[string]interface{}, prop string, expected string, ptr interface{}) error {
	var findErr error
	i := funk.Find(list, func(x map[string]interface{}) bool {
		if findErr != nil {
			return false
		}

		got := objx.New(x).Get(prop).Data()
		str, ok := got.(string)
		if !ok {
			findErr = fmt.Errorf("testing whether `%s` is `%s`: unable to parse `%v` as string", prop, expected, got)
			return false
		}

		return str == expected
	})
	if findErr != nil {
		return findErr
	}

	if i == nil {
		return ErrorNoMatch
	}

	o := objx.New(i).MustJSON()
	return json.Unmarshal([]byte(o), ptr)
}
