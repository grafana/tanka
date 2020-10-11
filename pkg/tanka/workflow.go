package tanka

import (
	"errors"
	"fmt"

	"github.com/fatih/color"

	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/term"
)

// ApplyOpts specify additional properties for the Apply action
type ApplyOpts struct {
	Opts

	// AutoApprove skips the interactive approval
	AutoApprove bool
	// DiffStrategy to use for printing the diff before approval
	DiffStrategy string

	// Force ignores any warnings kubectl might have
	Force bool
	// Validate set to false ignores invalid Kubernetes schemas
	Validate bool
}

// Apply parses the environment at the given directory (a `baseDir`) and applies
// the evaluated jsonnet to the Kubernetes cluster defined in the environments
// `spec.json`.
func Apply(path string, opts ApplyOpts) error {
	_, envs, err := ParseEnv(path, ParseOpts{JsonnetOpts: opts.JsonnetOpts})
	if err != nil {
		return err
	}

	for _, env := range envs {
		l, err := load(env, opts.Opts)
		if err != nil {
			return err
		}
		kube, err := l.connect()
		if err != nil {
			return err
		}
		defer kube.Close()

		// show diff
		diff, err := kube.Diff(l.Resources, kubernetes.DiffOpts{Strategy: opts.DiffStrategy})
		switch {
		case err != nil:
			// This is not fatal, the diff is not strictly required
			fmt.Println("Error diffing:", err)
		case diff == nil:
			tmp := "Warning: There are no differences. Your apply may not do anything at all."
			diff = &tmp
		}

		// in case of non-fatal error diff may be nil
		if diff != nil {
			b := term.Colordiff(*diff)
			fmt.Print(b.String())
		}

		// prompt for confirmation
		if opts.AutoApprove {
		} else if err := confirmPrompt("Applying to", l.Env.Spec.Namespace, kube.Info()); err != nil {
			if errors.Is(err, term.ErrConfirmationFailed) {
				fmt.Println(err)
				continue
			} else {
				return err
			}
		}

		if err = kube.Apply(l.Resources, kubernetes.ApplyOpts{
			Force:    opts.Force,
			Validate: opts.Validate,
		}); err != nil {
			return err
		}
	}
	return nil
}

// confirmPrompt asks the user for confirmation before apply
func confirmPrompt(action, namespace string, info client.Info) error {
	alert := color.New(color.FgRed, color.Bold).SprintFunc()

	return term.Confirm(
		fmt.Sprintf(`%s namespace '%s' of cluster '%s' at '%s' using context '%s'.`, action,
			alert(namespace),
			alert(info.Kubeconfig.Cluster.Name),
			alert(info.Kubeconfig.Cluster.Cluster.Server),
			alert(info.Kubeconfig.Context.Name),
		),
		"yes",
	)
}

// DiffOpts specify additional properties for the Diff action
type DiffOpts struct {
	Opts

	// Strategy must be one of "native" or "subset"
	Strategy string
	// Summarize prints a summary, instead of the actual diff
	Summarize bool
}

// Diff parses the environment at the given directory (a `baseDir`) and returns
// the differences from the live cluster state in `diff(1)` format. If the
// `WithDiffSummarize` modifier is used, a histogram created using `diffstat(1)`
// is returned instead.
// The cluster information is retrieved from the environments `spec.json`.
// NOTE: This function requires on `diff(1)`, `kubectl(1)` and perhaps `diffstat(1)`
func Diff(path string, opts DiffOpts) ([]*string, error) {
	_, envs, err := ParseEnv(path, ParseOpts{JsonnetOpts: opts.JsonnetOpts})
	if err != nil {
		return nil, err
	}

	diffs := make([]*string, 0)
	for _, env := range envs {
		l, err := load(env, opts.Opts)
		if err != nil {
			return nil, err
		}
		kube, err := l.connect()
		if err != nil {
			return nil, err
		}
		defer kube.Close()

		diff, err := kube.Diff(l.Resources, kubernetes.DiffOpts{
			Summarize: opts.Summarize,
			Strategy:  opts.Strategy,
		})
		if err != nil {
			return nil, err
		}
		if diff != nil {
			diffs = append(diffs, diff)
		}
	}
	if len(diffs) == 0 {
		return nil, nil
	}
	return diffs, nil
}

// DeleteOpts specify additional properties for the Delete operation
type DeleteOpts struct {
	Opts

	// AutoApprove skips the interactive approval
	AutoApprove bool

	// Force ignores any warnings kubectl might have
	Force bool
	// Validate set to false ignores invalid Kubernetes schemas
	Validate bool
}

// Delete parses the environment at the given directory (a `baseDir`) and deletes
// the generated objects from the Kubernetes cluster defined in the environment's
// `spec.json`.
func Delete(path string, opts DeleteOpts) error {
	_, envs, err := ParseEnv(path, ParseOpts{JsonnetOpts: opts.JsonnetOpts})
	if err != nil {
		return err
	}

	for _, env := range envs {
		l, err := load(env, opts.Opts)
		if err != nil {
			return err
		}
		kube, err := l.connect()
		if err != nil {
			return err
		}
		defer kube.Close()

		// show diff
		// static differ will never fail and always return something if input is not nil
		diff, err := kubernetes.StaticDiffer(false)(l.Resources)

		if err != nil {
			fmt.Println("Error diffing:", err)
		}

		// in case of non-fatal error diff may be nil
		if diff != nil {
			b := term.Colordiff(*diff)
			fmt.Print(b.String())
		}

		// prompt for confirmation
		if opts.AutoApprove {
		} else if err := confirmPrompt("Deleting from", l.Env.Spec.Namespace, kube.Info()); err != nil {
			if errors.Is(err, term.ErrConfirmationFailed) {
				fmt.Println(err)
				continue
			} else {
				return err
			}
		}

		if err = kube.Delete(l.Resources, kubernetes.DeleteOpts{
			Force:    opts.Force,
			Validate: opts.Validate,
		}); err != nil {
			return err
		}
	}
	return nil
}

// Show parses the environment at the given directory (a `baseDir`) and returns
// the list of Kubernetes objects.
// Tip: use the `String()` function on the returned list to get the familiar yaml stream
func Show(path string, opts Opts) (manifest.List, error) {
	_, envs, err := ParseEnv(path, ParseOpts{JsonnetOpts: opts.JsonnetOpts})
	if err != nil {
		return nil, err
	}
	resources := make(manifest.List, 0)
	for _, env := range envs {
		l, err := load(env, opts)
		if err != nil {
			return nil, err
		}
		resources = append(resources, l.Resources...)
	}
	return resources, nil
}

// Eval returns the raw evaluated Jsonnet output (without any transformations)
func Eval(dir string, opts Opts) (raw interface{}, err error) {
	r, _, err := ParseEnv(dir, ParseOpts{JsonnetOpts: opts.JsonnetOpts})
	switch err.(type) {
	case ErrNoEnv, ErrMultipleEnvs:
		return r, err
	case nil:
		return r, nil
	default:
		return nil, err
	}
}
