package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/posener/complete"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/grafana/tanka/pkg/cli/cmp"
	"github.com/grafana/tanka/pkg/kubernetes/util"
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
	force := cmd.Flags().Bool("force", false, "force applying (kubectl apply --force)")
	validate := cmd.Flags().Bool("validate", true, "validation of resources (kubectl --validate=false)")
	autoApprove := cmd.Flags().Bool("dangerous-auto-approve", false, "skip interactive approval. Only for automation!")
	getExtCode := extCodeParser(cmd.Flags())

	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := tanka.Apply(args[0],
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
			tanka.WithExtCode(getExtCode()),
			tanka.WithApplyForce(*force),
			tanka.WithApplyValidate(*validate),
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

	getExtCode := extCodeParser(cmd.Flags())

	cmd.Run = func(cmd *cobra.Command, args []string) {
		changes, err := tanka.Diff(args[0],
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
			tanka.WithExtCode(getExtCode()),
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
	getExtCode := extCodeParser(cmd.Flags())
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if !interactive && !*allowRedirect {
			fmt.Fprintln(os.Stderr, "Redirection of the output of tk show is discouraged and disabled by default. Run tk show --dangerous-allow-redirect to enable.")
			return
		}

		pretty, err := tanka.Show(args[0],
			tanka.WithExtCode(getExtCode()),
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
		)
		if err != nil {
			log.Fatalln(err)
		}

		pageln(pretty)
	}
	return cmd
}

func stringsToRegexps(exps []string) []*regexp.Regexp {
	regexs, err := util.CompileTargetExps(exps)
	if err != nil {
		log.Fatalln(err)
	}
	return regexs
}
