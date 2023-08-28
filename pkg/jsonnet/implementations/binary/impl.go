package binary

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/grafana/tanka/pkg/jsonnet/implementations/types"
)

type JsonnetBinaryRunner struct {
	binPath string
	args    []string
}

func (r *JsonnetBinaryRunner) EvaluateAnonymousSnippet(snippet string) (string, error) {
	cmd := exec.Command(r.binPath, append(r.args, "-e", snippet)...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running anonymous snippet: %w\n%s", err, string(out))
	}

	return string(out), nil
}

func (r *JsonnetBinaryRunner) EvaluateFile(filename string) (string, error) {
	cmd := exec.Command(r.binPath, append(r.args, filename)...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running file %s: %w\n%s", filename, err, string(out))
	}

	return string(out), nil
}

// JsonnetBinaryImplementation runs Jsonnet in a subprocess. It doesn't support native functions.
// The interface of the binary has to compatible with the official Jsonnet CLI.
// It has to support the following flags:
// -J <path> for specifying import paths
// --ext-code <name>=<value> for specifying external variables
// --tla-code <name>=<value> for specifying top-level arguments
// --max-stack <value> for specifying the maximum stack size
// -e <code> for evaluating code snippets
// <filename> positional arg for evaluating files
type JsonnetBinaryImplementation struct {
	BinPath string
}

func (i *JsonnetBinaryImplementation) MakeEvaluator(importPaths []string, extCode map[string]string, tlaCode map[string]string, maxStack int) types.JsonnetEvaluator {
	args := []string{}
	for _, p := range importPaths {
		args = append(args, "-J", p)
	}
	if maxStack > 0 {
		args = append(args, "--max-stack", strconv.Itoa(maxStack))
	}
	for k, v := range extCode {
		args = append(args, "--ext-code", k+"="+v)
	}
	for k, v := range tlaCode {
		args = append(args, "--tla-code", k+"="+v)
	}

	return &JsonnetBinaryRunner{
		binPath: i.BinPath,
		args:    args,
	}
}
