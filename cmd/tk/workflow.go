package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/posener/complete"
	"github.com/spf13/pflag"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/kubernetes/util"
	"github.com/grafana/tanka/pkg/tanka"
	"github.com/grafana/tanka/pkg/term"
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

func applyCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "apply <path>",
		Short: "apply the configuration to the cluster",
		Args:  workflowArgs,
	}

	vars := workflowFlags(cmd.Flags())
	force := cmd.Flags().Bool("force", false, "force applying (kubectl apply --force)")
	validate := cmd.Flags().Bool("validate", true, "validation of resources (kubectl --validate=false)")
	autoApprove := cmd.Flags().Bool("dangerous-auto-approve", false, "skip interactive approval. Only for automation!")
	getExtCode := extCodeParser(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		err := tanka.Apply(args[0],
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
			tanka.WithExtCode(getExtCode()),
			tanka.WithApplyForce(*force),
			tanka.WithApplyValidate(*validate),
			tanka.WithApplyAutoApprove(*autoApprove),
		)
		if err != nil {
			return err
		}
		return nil
	}
	return cmd
}

func diffCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "diff <path>",
		Short: "differences between the configuration and the cluster",
		Args:  workflowArgs,
		Predictors: complete.Flags{
			"diff-strategy": cli.PredictSet("native", "subset"),
		},
	}

	// flags
	var (
		vars         = workflowFlags(cmd.Flags())
		diffStrategy = cmd.Flags().String("diff-strategy", "", "force the diff-strategy to use. Automatically chosen if not set.")
		summarize    = cmd.Flags().BoolP("summarize", "s", false, "quick summary of the differences, hides file contents")
	)

	getExtCode := extCodeParser(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		changes, err := tanka.Diff(args[0],
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
			tanka.WithExtCode(getExtCode()),
			tanka.WithDiffStrategy(*diffStrategy),
			tanka.WithDiffSummarize(*summarize),
		)
		if err != nil {
			return err
		}

		if changes == nil {
			log.Println("No differences.")
			os.Exit(ExitStatusClean)
		}

		if interactive {
			r := term.Colordiff(*changes)
			fPageln(r)
		} else {
			fmt.Println(*changes)
		}

		os.Exit(ExitStatusDiff)
		return nil
	}

	return cmd
}

func showCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "show <path>",
		Short: "jsonnet as yaml",
		Args:  workflowArgs,
	}
	vars := workflowFlags(cmd.Flags())
	allowRedirect := cmd.Flags().Bool("dangerous-allow-redirect", false, "allow redirecting output to a file or a pipe.")
	getExtCode := extCodeParser(cmd.Flags())
	cmd.Run = func(cmd *cli.Command, args []string) error {
		if !interactive && !*allowRedirect {
			fmt.Fprintln(os.Stderr, `Redirection of the output of tk show is discouraged and disabled by default.
If you want to export .yaml files for use with other tools, try 'tk export'.
Otherwise run tk show --dangerous-allow-redirect to bypass this check.`)
			return nil
		}

		pretty, err := tanka.Show(args[0],
			tanka.WithExtCode(getExtCode()),
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
		)
		if err != nil {
			return err
		}

		pageln(pretty.String())
		return nil
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
