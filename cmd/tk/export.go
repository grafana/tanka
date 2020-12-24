package main

import (
	"fmt"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/tanka"
)

func exportCmd() *cli.Command {
	args := workflowArgs
	args.Validator = cli.ValidateFunc(func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("expects at least 2 args, received %v", len(args))
		}
		return nil
	})

	cmd := &cli.Command{
		Use:   "export <outputDir> <path> [<path>...]",
		Short: "export environments found in path(s)",
		Args:  args,
	}

	opts := tanka.DefaultExportEnvOpts()

	opts.Format = cmd.Flags().String("format", *opts.Format, "https://tanka.dev/exporting#filenames")

	opts.Extension = cmd.Flags().String("extension", *opts.Extension, "File extension")
	opts.Merge = cmd.Flags().Bool("merge", *opts.Merge, "Allow merging with existing directory")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())
	getLabelSelector := labelSelectorFlag(cmd.Flags())

	recursive := cmd.Flags().BoolP("recursive", "r", false, "Look recursively for Tanka environments")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		opts.Targets = vars.targets
		opts.ParseParallelOpts.ParseOpts.JsonnetOpts = getJsonnetOpts()
		opts.ParseParallelOpts.Selector = getLabelSelector()

		var paths []string
		for _, path := range args[1:] {
			if *recursive {
				envs, err := tanka.FindEnvironments(path, opts.ParseParallelOpts.Selector)
				if err != nil {
					return err
				}
				for _, env := range envs {
					paths = append(paths, env.Metadata.Namespace)
				}
			} else {
				parseOpts := tanka.ParseOpts{
					Evaluator: tanka.EnvsOnlyEvaluator,
				}
				_, _, err := tanka.ParseEnv(path, parseOpts)
				if err != nil {
					return err
				}
				paths = append(paths, path)
			}
		}

		return tanka.ExportEnvironments(paths, args[0], &opts)
	}
	return cmd
}
