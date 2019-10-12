package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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

func applyDeleteFlags(fs *pflag.FlagSet) kubernetes.ApplyDeleteOpts {
	force := fs.Bool("force", false, "force operating (add --force for kubelet)")
	autoApprove := fs.Bool("dangerous-auto-approve", false, "skip interactive approval. Only for automation!")

	return kubernetes.ApplyDeleteOpts{
		Force:       *force,
		AutoApprove: *autoApprove,
	}
}

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [path]",
		Short: "apply the configuration to the cluster",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	vars := workflowFlags(cmd.Flags())
	applyFlags := applyDeleteFlags(cmd.Flags())
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if kube == nil {
			log.Fatalln(kubernetes.ErrorMissingConfig{Verb: "apply"})
		}

		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		desired, err := kube.Reconcile(raw, stringsToRegexps(vars.targets)...)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		if !diff(desired, false, kubernetes.DiffOpts{}) {
			log.Println("Warning: There are no differences. Your apply may not do anything at all.")
		}

		if err := kube.Apply(desired, applyFlags); err != nil {
			log.Fatalln("Applying:", err)
		}
	}
	return cmd
}

func diffCmd() *cobra.Command {
	// completion
	cmp.Handlers.Add("diffStrategy", complete.PredictSet("native", "subset"))

	cmd := &cobra.Command{
		Use:   "diff [path]",
		Short: "differences between the configuration and the cluster",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args":                "baseDir",
			"flags/diff-strategy": "diffStrategy",
		},
	}

	// flags
	var (
		vars         = workflowFlags(cmd.Flags())
		diffStrategy = cmd.Flags().String("diff-strategy", "", "force the diff-strategy to use. Automatically chosen if not set.")
		summarize    = cmd.Flags().BoolP("summarize", "s", false, "quick summary of the differences, hides file contents")
	)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if kube == nil {
			log.Fatalln(kubernetes.ErrorMissingConfig{Verb: "diff"})
		}

		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		if kube.Spec.DiffStrategy == "" {
			kube.Spec.DiffStrategy = *diffStrategy
		}

		desired, err := kube.Reconcile(raw, stringsToRegexps(vars.targets)...)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		if diff(desired, interactive, kubernetes.DiffOpts{Summarize: *summarize}) {
			os.Exit(16)
		}
		log.Println("No differences.")
	}

	return cmd
}

// computes the diff, prints to screen.
// set `pager` to false to disable the pager.
// When non-interactive, no paging happens anyways.
func diff(state []kubernetes.Manifest, pager bool, opts kubernetes.DiffOpts) (changed bool) {
	changes, err := kube.Diff(state, opts)
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
		Use:   "show [path]",
		Short: "jsonnet as yaml",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	vars := workflowFlags(cmd.Flags())
	canRedirect := cmd.Flags().Bool("dangerous-allow-redirect", false, "allow redirecting output to a file or a pipe.")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if !interactive && !*canRedirect {
			fmt.Fprintln(os.Stderr, "Redirection of the output of tk show is discouraged and disabled by default. Run tk show --dangerous-allow-redirect to enable.")
			return
		}

		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		state, err := kube.Reconcile(raw, stringsToRegexps(vars.targets)...)
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

// stringsToRegexps compiles each string to a regular expression
func stringsToRegexps(strs []string) (exps []*regexp.Regexp) {
	exps = make([]*regexp.Regexp, 0, len(strs))
	for _, raw := range strs {
		// surround the regular expression with start and end markers
		s := fmt.Sprintf(`^%s$`, raw)
		exp, err := regexp.Compile(s)
		if err != nil {
			log.Fatalf("%s.\nSee https://tanka.dev/targets/#regular-expressions for details on regular expressions.", strings.Title(err.Error()))
		}
		exps = append(exps, exp)
	}
	return exps
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [path]",
		Short: "delete configuration from the cluster",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	vars := workflowFlags(cmd.Flags())
	deleteFlags := applyDeleteFlags(cmd.Flags())
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if kube == nil {
			log.Fatalln(kubernetes.ErrorMissingConfig{Verb: "delete"})
		}

		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		desired, err := kube.Reconcile(raw, stringsToRegexps(vars.targets)...)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		if err := kube.Delete(desired, deleteFlags); err != nil {
			log.Fatalln("Deleting:", err)
		}
	}
	return cmd
}
