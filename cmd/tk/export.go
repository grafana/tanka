package main

import (
	"fmt"
	"path/filepath"

	"github.com/go-clix/cli"
	"github.com/pkg/errors"

	"github.com/grafana/tanka/pkg/export"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/process"
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

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())
	getLabelSelector := labelSelectorFlag(cmd.Flags())
	recursive := cmd.Flags().BoolP("recursive", "r", false, "Look recursively for Tanka environments")

	var opts export.Opts
	cmd.Flags().StringVar(&opts.Extension, "extension", "yaml", "File extension")
	cmd.Flags().BoolVar(&opts.Merge, "merge", false, "Allow merging with existing directory")
	cmd.Flags().StringVar(&opts.Format, "format", export.DefaultFormat, "https://tanka.dev/exporting#filenames")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		opts.Filters = process.MustStrExps(vars.targets...) // TODO: check err
		opts.JsonnetOpts = getJsonnetOpts()

		outputDir := args[0]

		var paths []string
		for _, path := range args[1:] {
			// recursive?
			if *recursive {
				rootDir, err := jpath.FindRoot(path)
				if err != nil {
					return errors.Wrap(err, "resolving jpath")
				}

				envs, err := tanka.ListEnvs(path, tanka.ListOpts{
					Selector: getLabelSelector(),
				})
				if err != nil {
					return err
				}

				for _, env := range envs {
					paths = append(paths, filepath.Join(rootDir, env.Metadata.Namespace))
				}

				continue
			}

			paths = append(paths, path)
		}

		// export them
		return export.Export(paths, outputDir, &opts)
	}
	return cmd
}
