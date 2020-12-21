package main

import (
	"github.com/go-clix/cli"
	"github.com/grafana/tanka/pkg/tanka"
)

func exportCmd() *cli.Command {
	args := workflowArgs
	args.Validator = cli.ValidateExact(2)

	cmd := &cli.Command{
		Use:   "export <environment> <outputDir>",
		Short: "write each resources as a YAML file",
		Args:  args,
	}

	defaultOpts := tanka.DefaultExportEnvOpts()

	format := cmd.Flags().String("format", defaultOpts.Format, "https://tanka.dev/exporting#filenames")
	dirFormat := cmd.Flags().String("dirformat", defaultOpts.DirFormat, "based on tanka.dev/Environment object")

	extension := cmd.Flags().String("extension", defaultOpts.Extension, "File extension")
	merge := cmd.Flags().Bool("merge", defaultOpts.Merge, "Allow merging with existing directory")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		opts := tanka.ExportEnvOpts{
			Format:      *format,
			DirFormat:   *dirFormat,
			Extension:   *extension,
			Targets:     vars.targets,
			Merge:       *merge,
			JsonnetOpts: getJsonnetOpts(),
		}
		return tanka.ExportEnvironments([]string{args[0]}, args[1], &opts)
	}
	return cmd
}
