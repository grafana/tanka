package main

import (
	"fmt"
	"path/filepath"

	"github.com/go-clix/cli"
	"github.com/pkg/errors"

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
		opts := tanka.ExportOpts{
			Format:    *format,
			Extension: *extension,
			Merge:     *merge,
			Opts: tanka.Opts{
				Filters:     process.MustStrExps(vars.targets...), // TODO: check err
				JsonnetOpts: getJsonnetOpts(),
			},
		}

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

			// single env
			paths = []string{path}
		}

		// export them
		return tanka.Export(paths, outputDir, &opts)
	}
	return cmd
}
