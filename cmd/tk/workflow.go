package main

import (
	"fmt"
	"log"
	"os"

	"github.com/posener/complete"
	"github.com/spf13/pflag"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/process"
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

	var opts tanka.ApplyOpts
	cmd.Flags().BoolVar(&opts.Force, "force", false, "force applying (kubectl apply --force)")
	cmd.Flags().BoolVar(&opts.Validate, "validate", true, "validation of resources (kubectl --validate=false)")
	cmd.Flags().BoolVar(&opts.AutoApprove, "dangerous-auto-approve", false, "skip interactive approval. Only for automation!")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		opts.Filters = stringsToRegexps(vars.targets)
		opts.JsonnetOpts = getJsonnetOpts()

		return tanka.Apply(args[0], opts)
	}
	return cmd
}

func pruneCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "prune <path>",
		Short: "delete resources removed from Jsonnet",
		Args:  workflowArgs,
	}

	var opts tanka.PruneOpts
	cmd.Flags().BoolVar(&opts.Force, "force", false, "force deleting (kubectl delete --force)")
	cmd.Flags().BoolVar(&opts.AutoApprove, "dangerous-auto-approve", false, "skip interactive approval. Only for automation!")
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		opts.JsonnetOpts = getJsonnetOpts()

		return tanka.Prune(args[0], opts)
	}

	return cmd
}

func deleteCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "delete <path>",
		Short: "delete the environment from cluster",
		Args:  workflowArgs,
	}

	var opts tanka.DeleteOpts
	cmd.Flags().BoolVar(&opts.Force, "force", false, "force deleting (kubectl delete --force)")
	cmd.Flags().BoolVar(&opts.Validate, "validate", true, "validation of resources (kubectl --validate=false)")
	cmd.Flags().BoolVar(&opts.AutoApprove, "dangerous-auto-approve", false, "skip interactive approval. Only for automation!")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		opts.Filters = stringsToRegexps(vars.targets)
		opts.JsonnetOpts = getJsonnetOpts()

		return tanka.Delete(args[0], opts)
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

	var opts tanka.DiffOpts
	cmd.Flags().StringVar(&opts.Strategy, "diff-strategy", "", "force the diff-strategy to use. Automatically chosen if not set.")
	cmd.Flags().BoolVarP(&opts.Summarize, "summarize", "s", false, "print summary of the differences, not the actual contents")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		opts.Filters = stringsToRegexps(vars.targets)
		opts.JsonnetOpts = getJsonnetOpts()

		changes, err := tanka.Diff(args[0], opts)
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

	allowRedirect := cmd.Flags().Bool("dangerous-allow-redirect", false, "allow redirecting output to a file or a pipe.")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if !interactive && !*allowRedirect {
			fmt.Fprintln(os.Stderr, `Redirection of the output of tk show is discouraged and disabled by default.
If you want to export .yaml files for use with other tools, try 'tk export'.
Otherwise run tk show --dangerous-allow-redirect to bypass this check.`)
			return nil
		}

		pretty, err := tanka.Show(args[0], tanka.Opts{
			JsonnetOpts: getJsonnetOpts(),
			Filters:     stringsToRegexps(vars.targets),
		})

		if err != nil {
			return err
		}

		pageln(pretty.String())
		return nil
	}
	return cmd
}

func stringsToRegexps(exps []string) process.Matchers {
	regexs, err := process.StrExps(exps...)
	if err != nil {
		log.Fatalln(err)
	}
	return regexs
}
