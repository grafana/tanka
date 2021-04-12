package main

import (
	"fmt"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/process"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
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
		"{{.apiVersion}}.{{.kind}}-{{or .metadata.name .metadata.generateName}}",
		"https://tanka.dev/exporting#filenames",
	)

	extension := cmd.Flags().String("extension", "yaml", "File extension")
	merge := cmd.Flags().Bool("merge", false, "Allow merging with existing directory")
	parallel := cmd.Flags().IntP("parallel", "p", 8, "Number of environments to process in parallel")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())
	getLabelSelector := labelSelectorFlag(cmd.Flags())

	recursive := cmd.Flags().BoolP("recursive", "r", false, "Look recursively for Tanka environments")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		filters, err := process.StrExps(vars.targets...)
		if err != nil {
			return err
		}

		opts := tanka.ExportEnvOpts{
			Format:    *format,
			Extension: *extension,
			Merge:     *merge,
			Opts: tanka.Opts{
				JsonnetOpts: getJsonnetOpts(),
				Filters:     filters,
				Name:        vars.name,
			},
			Selector:    getLabelSelector(),
			Parallelism: *parallel,
		}

		var exportEnvs []*v1alpha1.Environment
		for _, path := range args[1:] {
			// find possible environments
			if *recursive {
				// get absolute path to Environment
				envs, err := tanka.FindEnvs(path, tanka.FindOpts{Selector: opts.Selector})
				if err != nil {
					return err
				}

				for _, env := range envs {
					if opts.Opts.Name != "" && opts.Opts.Name != env.Metadata.Name {
						continue
					}
					exportEnvs = append(exportEnvs, env)
				}
				continue
			}

			// validate environment
			env, err := tanka.Peek(path, opts.Opts)
			if err != nil {
				switch err.(type) {
				case tanka.ErrMultipleEnvs:
					fmt.Println("Please use --name to export a single environment or --recursive to export multiple environments.")
					return err
				default:
					return err
				}
			}

			exportEnvs = append(exportEnvs, env)
		}

		// export them
		return tanka.ExportEnvironments(exportEnvs, args[0], &opts)
	}
	return cmd
}
