package main

import (
	"fmt"
	"log"
	"os"

	"github.com/posener/complete"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/grafana/tanka/pkg/cmp"
	"github.com/grafana/tanka/pkg/kubernetes"
)

type workflowFlagVars struct {
	targets []string
}

func workflowFlags(fs *pflag.FlagSet) *workflowFlagVars {
	v := workflowFlagVars{}
	fs.StringSliceVarP(&v.targets, "target", "t", nil, "only use the specified objects (Format: <type>/<name>)")
	return &v
}

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [directory]",
		Short: "apply the configuration to the cluster",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	vars := workflowFlags(cmd.Flags())
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		desired, err := kube.Reconcile(raw, vars.targets...)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		if !diff(desired, false) {
			log.Println("Warning: There are no differences. Your apply may not do anything at all.")
		}

		if err := kube.Apply(desired); err != nil {
			log.Fatalln("Applying:", err)
		}
	}
	return cmd
}

func diffCmd() *cobra.Command {
	// completion
	cmp.Handlers.Add("diffStrategy", complete.PredictSet("native", "subset"))

	cmd := &cobra.Command{
		Use:   "diff [directory]",
		Short: "differences between the configuration and the cluster",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args":                "baseDir",
			"flags/diff-strategy": "diffStrategy",
		},
	}
	vars := workflowFlags(cmd.Flags())
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		if kube.Spec.DiffStrategy == "" {
			kube.Spec.DiffStrategy = cmd.Flag("diff-strategy").Value.String()
		}

		desired, err := kube.Reconcile(raw, vars.targets...)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		if diff(desired, interactive) {
			os.Exit(16)
		}
		log.Println("No differences.")
	}

	cmd.Flags().String("diff-strategy", "", "force the diff-strategy to use. Automatically chosen if not set.")
	return cmd
}

// computes the diff, prints to screen.
// set `pager` to false to disable the pager.
// When non-interactive, no paging happens anyways.
func diff(state []kubernetes.Manifest, pager bool) (changed bool) {
	changes, err := kube.Diff(state)
	if err != nil {
		log.Fatalln("Diffing:", err)
	}

	if changes == nil {
		return false
	}

	if interactive {
		h := highlight("diff", *changes)
		if pager {
			pageln(h)
		} else {
			fmt.Println(h)
		}
	} else {
		fmt.Println(*changes)
	}

	return true
}

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [directory]",
		Short: "jsonnet as yaml",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	vars := workflowFlags(cmd.Flags())
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		state, err := kube.Reconcile(raw, vars.targets...)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		pretty, err := kube.Fmt(state)
		if err != nil {
			log.Fatalln("Pretty printing state:", err)
		}

		pageln(pretty)
	}
	return cmd
}
