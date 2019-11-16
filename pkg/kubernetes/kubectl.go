package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/stretchr/objx"
	funk "github.com/thoas/go-funk"

	"github.com/grafana/tanka/pkg/cli"
)

var (
	alert = color.New(color.FgRed, color.Bold).SprintFunc()
)

// Kubectl uses the `kubectl` command to operate on a Kubernetes cluster
type Kubectl struct {
	context   objx.Map
	cluster   objx.Map
	APIServer string
}

// Version returns the version of kubectl and the Kubernetes api server
func (k Kubectl) Version() (client, server semver.Version, err error) {
	zero := *semver.MustParse("0.0.0")
	if err := k.setupContext(); err != nil {
		return zero, zero, err
	}
	cmd := exec.Command("kubectl", "version",
		"-o", "json",
		"--context", k.context.Get("name").MustStr(),
	)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return zero, zero, err
	}
	vs := objx.MustFromJSON(buf.String())
	client = *semver.MustParse(vs.Get("clientVersion.gitVersion").MustStr())
	server = *semver.MustParse(vs.Get("serverVersion.gitVersion").MustStr())
	return client, server, nil
}

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

// Get retrieves an Kubernetes object from the API
func (k Kubectl) Get(namespace, kind, name string) (map[string]interface{}, error) {
	if err := k.setupContext(); err != nil {
		return nil, err
	}
	argv := []string{"get",
		"-o", "json",
		"-n", namespace,
		"--context", k.context.Get("name").MustStr(),
		kind, name,
	}
	cmd := exec.Command("kubectl", argv...)
	var sout, serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	if err := cmd.Run(); err != nil {
		if strings.HasPrefix(serr.String(), "Error from server (NotFound)") {
			return nil, ErrorNotFound{kind, name}
		}
		fmt.Print(serr.String())
		return nil, err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(sout.Bytes(), &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// ApplyDeleteOpts allow to specify additional parameter for apply/delete operations
type ApplyDeleteOpts struct {
	Force       bool
	AutoApprove bool
}

// Apply applies the given yaml to the cluster
func (k Kubectl) Apply(yaml, namespace string, opts ApplyDeleteOpts) error {
	if err := k.setupContext(); err != nil {
		return err
	}

	if !opts.AutoApprove {
		if err := cli.Confirm(
			fmt.Sprintf(`Applying to namespace '%s' of cluster '%s' at '%s' using context '%s'.`,
				alert(namespace),
				alert(k.cluster.Get("name").MustStr()),
				alert(k.cluster.Get("cluster.server").MustStr()),
				alert(k.context.Get("name").MustStr()),
			),
			"yes",
		); err != nil {
			return err
		}
	}

	argv := []string{"apply",
		"--context", k.context.Get("name").MustStr(),
		"-f", "-",
	}
	if !opts.Force {
		argv = append(argv, "--force")
	}

	cmd := exec.Command("kubectl", argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		fmt.Fprintln(stdin, yaml)
		stdin.Close()
	}()

	return cmd.Run()
}

// Diff takes a desired state as yaml and returns the differences
// to the system in common diff format
func (k Kubectl) Diff(yaml string) (*string, error) {
	if err := k.setupContext(); err != nil {
		return nil, err
	}

	argv := []string{"diff",
		"--context", k.context.Get("name").MustStr(),
		"-f", "-",
	}
	cmd := exec.Command("kubectl", argv...)
	raw := bytes.Buffer{}
	cmd.Stdout = &raw
	cmd.Stderr = FilteredErr{regexp.MustCompile(`exit status \d`)}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		fmt.Fprintln(stdin, yaml)
		stdin.Close()
	}()

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// kubectl uses this to tell us that there is a diff
			if exitError.ExitCode() == 1 {
				s := raw.String()
				return &s, nil
			}
		}
		return nil, err
	}

	// no diff -> nil
	return nil, nil
}

// Delete removes the given yaml from the cluster
func (k Kubectl) Delete(yaml, namespace string, opts ApplyDeleteOpts) error {
	if err := k.setupContext(); err != nil {
		return err
	}

	if !opts.AutoApprove {
		if err := util.Confirm(
			fmt.Sprintf(`Deleting from namespace '%s' of cluster '%s' at '%s' using context '%s'.`,
				alert(namespace),
				alert(k.cluster.Get("name").MustStr()),
				alert(k.cluster.Get("cluster.server").MustStr()),
				alert(k.context.Get("name").MustStr()),
			),
			"yes",
		); err != nil {
			return err
		}
	}

	argv := []string{"delete",
		"--context", k.context.Get("name").MustStr(),
		"-f", "-",
	}
	if !opts.Force {
		argv = append(argv, "--force")
	}

	cmd := exec.Command("kubectl", argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		fmt.Fprintln(stdin, yaml)
		stdin.Close()
	}()

	return cmd.Run()
}
