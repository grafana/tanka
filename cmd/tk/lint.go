package main

import (
	"errors"

	"github.com/go-clix/cli"
	"github.com/gobwas/glob"
	"github.com/posener/complete"

	"github.com/grafana/tanka/pkg/jsonnet"
)

func lintCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "lint <FILES|DIRECTORIES>",
		Short: "lint Jsonnet code",
		Args: cli.Args{
			Validator: cli.ValidateFunc(func(args []string) error {
				if len(args) == 0 {
					return errors.New("at least one file or directory is required")
				}
				return nil
			}),
			Predictor: complete.PredictFiles("*.*sonnet"),
		},
	}

	exclude := cmd.Flags().StringSliceP("exclude", "e", []string{"**/.*", ".*", "**/vendor/**", "vendor/**"}, "globs to exclude")
	parallelism := cmd.Flags().IntP("parallelism", "n", 4, "amount of workers")
	verbose := cmd.Flags().BoolP("verbose", "v", false, "print each checked file")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		globs := make([]glob.Glob, len(*exclude))
		for i, e := range *exclude {
			g, err := glob.Compile(e)
			if err != nil {
				return err
			}
			globs[i] = g
		}

		return jsonnet.Lint(args, &jsonnet.LintOpts{Excludes: globs, PrintNames: *verbose, Parallelism: *parallelism})
	}

	return cmd
}
