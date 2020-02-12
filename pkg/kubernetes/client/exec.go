package client

import (
	"fmt"
	"os"
	"os/exec"
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

	// prepend namespace patch to $KUBECONFIG
	env := newEnv()
	if _, ok := env["KUBECONFIG"]; !ok {
		env["KUBECONFIG"] = "~/.kube/config"
	}
	env["KUBECONFIG"] = fmt.Sprintf("%s:%s", k.nsPatch, env["KUBECONFIG"])
	cmd.Env = env.render()

	return cmd
}

// environment is a helper type for manipulating os.Environ() more easily
type environment map[string]string

func newEnv() environment {
	env := make(environment)
	for _, s := range os.Environ() {
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
	return s
}
