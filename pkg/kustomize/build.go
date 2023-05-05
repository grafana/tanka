package kustomize

import (
	"bytes"
	"io"
	"os"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

func (k ExecKustomize) buildCommandArgs(path string, opts BuildOpts) []string {
	args := []string{path}
	args = append(args, opts.Flags()...)
	return args
}

// Build expands a Kustomize into a regular manifest.List using the `kustomize
// build` command
func (k ExecKustomize) Build(path string, opts BuildOpts) (manifest.List, error) {
	args := k.buildCommandArgs(path, opts)

	cmd := k.cmd("build", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "Expanding Kustomize")
	}

	var list manifest.List
	d := yaml.NewDecoder(&buf)
	for {
		var m manifest.Manifest
		if err := d.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "Parsing Kustomize output")
		}

		list = append(list, m)
	}

	return list, nil
}

// BuildOpts are additional, non-required options for Kustomize.Build
type BuildOpts struct {
	// enable kustomize plugins
	EnableAlphaPlugins bool
	// enable support for exec functions (raw executables); do not use for untrusted configs! (Alpha)
	EnableExec bool
	// Enable use of the Helm chart inflator generator.
	EnableHelm bool
	// enable support for starlark functions. (Alpha)
	EnableStar bool
	// if set to 'LoadRestrictionsNone', local kustomizations may load files from outside their root.
	// This does, however, break the relocatability of the kustomization. (default "LoadRestrictionsRootOnly")
	LoadRestrictor string
}

// Flags returns all options as their respective `kustomize build` flag equivalent
func (b BuildOpts) Flags() []string {
	var flags []string

	if b.EnableAlphaPlugins {
		flags = append(flags, "--enable-alpha-plugins")
	}

	if b.EnableExec {
		flags = append(flags, "--enable-exec")
	}

	if b.EnableHelm {
		flags = append(flags, "--enable-helm")
	}

	if b.EnableStar {
		flags = append(flags, "--enable-star")
	}

	if b.LoadRestrictor != "" {
		flags = append(flags, "--load-restrictor="+b.LoadRestrictor)
	}

	return flags
}
