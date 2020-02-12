package client

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// ctl returns an `exec.Cmd` for `kubectl`. It also forces the correct context
// and injects our patched $KUBECONFIG for the default namespace.
func (k Kubectl) ctl(action string, args ...string) *exec.Cmd {
	// prepare the arguments
	argv := []string{action,
		"--context", k.context.Get("name").MustStr(),
	}
	argv = append(argv, args...)

	// prepare the cmd
	cmd := exec.Command("kubectl", argv...)
	cmd.Env = patchKubeconfig(k.nsPatch, os.Environ())

	return cmd
}

func patchKubeconfig(file string, e []string) []string {
	// prepend namespace patch to $KUBECONFIG
	env := newEnv(e)
	if _, ok := env["KUBECONFIG"]; !ok {
		env["KUBECONFIG"] = "~/.kube/config" // kubectl default
	}
	env["KUBECONFIG"] = fmt.Sprintf("%s:%s", file, env["KUBECONFIG"])
	return env.render()
}

// environment is a helper type for manipulating os.Environ() more easily
type environment map[string]string

func newEnv(e []string) environment {
	env := make(environment)
	for _, s := range e {
		kv := strings.SplitN(s, "=", 2)
		env[kv[0]] = kv[1]
	}
	return env
}

func (e environment) render() []string {
	s := make([]string, 0, len(e))
	for k, v := range e {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(s)
	return s
}
