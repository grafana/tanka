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

	"github.com/grafana/tanka/pkg/cli/cmp"
	"github.com/grafana/tanka/pkg/tanka"
)

// special exit codes for tk diff
const (
	// no changes
	ExitStatusClean = 0
	// differences between the local config and the cluster
	ExitStatusDiff = 16
)

type workflowFlagVars struct {
	targets []string
}

func workflowFlags(fs *pflag.FlagSet) *workflowFlagVars {
	v := workflowFlagVars{}
	fs.StringSliceVarP(&v.targets, "target", "t", nil, "only use the specified objects (Format: <type>/<name>)")
	return &v
}

func applyDeleteFlags(fs *pflag.FlagSet) *kubernetes.ApplyDeleteOpts {
	opts := kubernetes.ApplyDeleteOpts{}
	fs.BoolVar(&opts.Force, "force", false, "force operating (kubectl --force)")
	fs.BoolVar(&opts.AutoApprove, "dangerous-auto-approve", false, "skip interactive approval. Only for automation!")
	return &opts
}

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply <path>",
		Short: "apply the configuration to the cluster",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	vars := workflowFlags(cmd.Flags())
	applyFlags := applyDeleteFlags(cmd.Flags())
	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := tanka.Apply(args[0],
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
			tanka.WithApplyForce(*force),
			tanka.WithApplyAutoApprove(*autoApprove),
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return cmd
}

func diffCmd() *cobra.Command {
	// completion
	cmp.Handlers.Add("diffStrategy", complete.PredictSet("native", "subset"))

	cmd := &cobra.Command{
		Use:   "diff <path>",
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
		changes, err := tanka.Diff(args[0],
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
			tanka.WithDiffStrategy(*diffStrategy),
			tanka.WithDiffSummarize(*summarize),
		)
		if err != nil {
			log.Fatalln(err)
		}

		if changes == nil {
			log.Println("No differences.")
			os.Exit(ExitStatusClean)
		}

		if interactive {
			h := highlight("diff", *changes)
			pageln(h)
		} else {
			fmt.Println(*changes)
		}

		os.Exit(ExitStatusDiff)
	}

	return cmd
}

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <path>",
		Short: "jsonnet as yaml",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	vars := workflowFlags(cmd.Flags())
	allowRedirect := cmd.Flags().Bool("dangerous-allow-redirect", false, "allow redirecting output to a file or a pipe.")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if !interactive && !*allowRedirect {
			fmt.Fprintln(os.Stderr, "Redirection of the output of tk show is discouraged and disabled by default. Run tk show --dangerous-allow-redirect to enable.")
			return
		}

		pretty, err := tanka.Show(args[0],
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
		)
		if err != nil {
			log.Fatalln(err)
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
		Use:   "delete <path>",
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

		if err := kube.Delete(desired, *deleteFlags); err != nil {
			log.Fatalln("Deleting:", err)
		}
	}
	return cmd
}
