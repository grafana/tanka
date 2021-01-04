package main

import (
	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/tanka"
)

func exportCmd() *cli.Command {
	args := workflowArgs
	args.Validator = ValidateMin(2)

	cmd := &cli.Command{
		Use:   "export <outputDir> <path> [<path>...]",
		Short: "export environments found in path(s)",
		Args:  args,
	}

	format := cmd.Flags().String(
		"format",
		"{{.apiVersion}}.{{.kind}}-{{.metadata.name}}",
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
			ParallelOpts: tanka.ParallelOpts{
				JsonnetOpts: getJsonnetOpts(),
				Selector:    getLabelSelector(),
			},
		}

		var paths []string
		for _, path := range args[1:] {
			// find possible environments
			if *recursive {
				// get absolute path to Environment
				envs, err := tanka.FindEnvironments(path, opts.ParallelOpts.Selector)
				if err != nil {
					return err
				}

				for path := range envs {
					paths = append(paths, path)
				}
				continue
			}

			// validate environment
			jsonnetOpts := opts.ParallelOpts.JsonnetOpts
			jsonnetOpts.EvalScript = tanka.EnvsOnlyEvalScript
			_, err := tanka.Load(path, tanka.Opts{JsonnetOpts: jsonnetOpts})
			if err != nil {
				return err
			}
			paths = append(paths, path)
		}

		// export them
		return tanka.ExportEnvironments(paths, args[0], &opts)
	}
	return cmd
}
