package main

import (
	"fmt"
	"path/filepath"

	"github.com/go-clix/cli"
	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
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
			ParseParallelOpts: tanka.ParseParallelOpts{
				JsonnetOpts: getJsonnetOpts(),
				Selector:    getLabelSelector(),
			},
		}

		var paths []string
		for _, path := range args[1:] {
			// find possible environments
			if *recursive {
				rootDir, err := jpath.FindRoot(path)
				if err != nil {
					return errors.Wrap(err, "resolving jpath")
				}

				// get absolute path to Environment
				envs, err := tanka.FindEnvironments(path, opts.ParseParallelOpts.Selector)
				if err != nil {
					return err
				}

				for _, env := range envs {
					paths = append(paths, filepath.Join(rootDir, env.Metadata.Namespace))
				}
				continue
			}

			// validate environment
			if _, err := tanka.Peek(path, opts.ParseParallelOpts.JsonnetOpts); err != nil {
				return err
			}

			paths = append(paths, path)
		}

		// export them
		return tanka.ExportEnvironments(paths, args[0], &opts)
	}
	return cmd
}
