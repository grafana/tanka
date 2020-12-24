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

	format := cmd.Flags().String(
		"format",
		"{{env.spec.namespace}}/{{env.metadata.name}}/{{.apiVersion}}.{{.kind}}-{{.metadata.name}}",
		"https://tanka.dev/exporting#filenames",
	)

	extension := cmd.Flags().String("extension", "yaml", "File extension")
	merge := cmd.Flags().Bool("merge", false, "Allow merging with existing directory")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())
	getLabelSelector := labelSelectorFlag(cmd.Flags())

	recursive := cmd.Flags().BoolP("recursive", "r", false, "Look recursively for Tanka environments")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		opts := tanka.ExportEnvOpts{
			Format:    *format,
			Extension: *extension,
			Merge:     *merge,
			Targets:   vars.targets,
			ParseParallelOpts: tanka.ParseParallelOpts{
				ParseOpts: tanka.ParseOpts{
					JsonnetOpts: getJsonnetOpts(),
				},
				Selector: getLabelSelector(),
			},
		}

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
				parseOpts := opts.ParseParallelOpts.ParseOpts
				parseOpts.Evaluator = tanka.EnvsOnlyEvaluator
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
