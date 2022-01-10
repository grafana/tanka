package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-clix/cli"
	"github.com/posener/complete"

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

func validateDryRun(dryRunStr string) error {
	switch dryRunStr {
	case "", "none", "client", "server":
		return nil
	}
	return fmt.Errorf(`--dry-run must be either: "", "none", "server" or "client"`)
}

func applyCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "apply <path>",
		Short: "apply the configuration to the cluster",
		Args:  workflowArgs,
		Predictors: complete.Flags{
			"diff-strategy":  cli.PredictSet("native", "subset", "validate", "server"),
			"apply-strategy": cli.PredictSet("client", "server"),
		},
	}

	var opts tanka.ApplyOpts
	cmd.Flags().BoolVar(&opts.Force, "force", false, "force applying (kubectl apply --force)")
	cmd.Flags().BoolVar(&opts.Validate, "validate", true, "validation of resources (kubectl --validate=false)")
	cmd.Flags().BoolVar(&opts.AutoApprove, "dangerous-auto-approve", false, "skip interactive approval. Only for automation!")
	cmd.Flags().StringVar(&opts.DryRun, "dry-run", "", `--dry-run parameter to pass down to kubectl, must be "none", "server", or "client"`)
	cmd.Flags().StringVar(&opts.ApplyStrategy, "apply-strategy", "", "force the apply strategy to use. Automatically chosen if not set.")
	cmd.Flags().StringVar(&opts.DiffStrategy, "diff-strategy", "", "force the diff strategy to use. Automatically chosen if not set.")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		err := validateDryRun(opts.DryRun)
		if err != nil {
			return err
		}

		filters, err := process.StrExps(vars.targets...)
		if err != nil {
			return err
		}
		opts.Filters = filters
		opts.JsonnetOpts = getJsonnetOpts()
		opts.Name = vars.name

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
	cmd.Flags().StringVar(&opts.DryRun, "dry-run", "", `--dry-run parameter to pass down to kubectl, must be "none", "server", or "client"`)
	cmd.Flags().BoolVar(&opts.Force, "force", false, "force deleting (kubectl delete --force)")
	cmd.Flags().BoolVar(&opts.AutoApprove, "dangerous-auto-approve", false, "skip interactive approval. Only for automation!")
	cmd.Flags().StringVar(&opts.Name, "name", "", "string that only a single inline environment contains in its name")
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		err := validateDryRun(opts.DryRun)
		if err != nil {
			return err
		}

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
	cmd.Flags().StringVar(&opts.DryRun, "dry-run", "", `--dry-run parameter to pass down to kubectl, must be "none", "server", or "client"`)
	cmd.Flags().BoolVar(&opts.Force, "force", false, "force deleting (kubectl delete --force)")
	cmd.Flags().BoolVar(&opts.Validate, "validate", true, "validation of resources (kubectl --validate=false)")
	cmd.Flags().BoolVar(&opts.AutoApprove, "dangerous-auto-approve", false, "skip interactive approval. Only for automation!")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		err := validateDryRun(opts.DryRun)
		if err != nil {
			return err
		}

		filters, err := process.StrExps(vars.targets...)
		if err != nil {
			return err
		}
		opts.Filters = filters
		opts.JsonnetOpts = getJsonnetOpts()
		opts.Name = vars.name

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
			"diff-strategy": cli.PredictSet("native", "subset", "validate", "server"),
		},
	}

	var opts tanka.DiffOpts
	cmd.Flags().StringVar(&opts.Strategy, "diff-strategy", "", "force the diff-strategy to use. Automatically chosen if not set.")
	cmd.Flags().BoolVarP(&opts.Summarize, "summarize", "s", false, "print summary of the differences, not the actual contents")
	cmd.Flags().BoolVarP(&opts.WithPrune, "with-prune", "p", false, "include objects deleted from the configuration in the differences")
	cmd.Flags().BoolVarP(&opts.ExitZero, "exit-zero", "z", false, "Exit with 0 even when differences are found.")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		filters, err := process.StrExps(vars.targets...)
		if err != nil {
			return err
		}
		opts.Filters = filters
		opts.JsonnetOpts = getJsonnetOpts()
		opts.Name = vars.name

		changes, err := tanka.Diff(args[0], opts)
		if err != nil {
			return err
		}

		if changes == nil {
			log.Println("No differences.")
			os.Exit(ExitStatusClean)
		}

		r := term.Colordiff(*changes)
		if err := fPageln(r); err != nil {
			return err
		}

		exitStatusDiff := ExitStatusDiff
		if opts.ExitZero {
			exitStatusDiff = ExitStatusClean
		}
		os.Exit(exitStatusDiff)
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

		filters, err := process.StrExps(vars.targets...)
		if err != nil {
			return err
		}

		pretty, err := tanka.Show(args[0], tanka.Opts{
			JsonnetOpts: getJsonnetOpts(),
			Filters:     filters,
			Name:        vars.name,
		})

		if err != nil {
			return err
		}

		return pageln(pretty.String())
	}
	return cmd
}
