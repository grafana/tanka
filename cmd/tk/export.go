package main

import (
	"fmt"

	"github.com/go-clix/cli"
	"k8s.io/apimachinery/pkg/labels"

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

	defaultOpts := tanka.DefaultExportEnvOpts()

	format := cmd.Flags().String("format", defaultOpts.Format, "https://tanka.dev/exporting#filenames")
	dirFormat := cmd.Flags().String("dirformat", defaultOpts.DirFormat, "based on tanka.dev/Environment object")

	extension := cmd.Flags().String("extension", defaultOpts.Extension, "File extension")
	merge := cmd.Flags().Bool("merge", defaultOpts.Merge, "Allow merging with existing directory")
	labelSelector := cmd.Flags().StringP("selector", "l", "", "Label selector. Uses the same syntax as kubectl does")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		var selector labels.Selector
		var err error
		if *labelSelector != "" {
			selector, err = labels.Parse(*labelSelector)
			if err != nil {
				return err
			}
		}

		opts := tanka.ExportEnvOpts{
			Format:    *format,
			DirFormat: *dirFormat,
			Extension: *extension,
			Targets:   vars.targets,
			Merge:     *merge,
			ParseOpts: tanka.ParseOpts{
				Selector:    selector,
				JsonnetOpts: getJsonnetOpts(),
			},
		}
		var paths []string
		for _, path := range args[1:] {
			dirs, err := tanka.FindBaseDirs(path)
			if err != nil {
				return err
			}
			paths = append(paths, dirs...)
		}
		return tanka.ExportEnvironments(paths, args[0], &opts)
	}
	return cmd
}
