package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-clix/cli"
	"github.com/gobwas/glob"
	"github.com/posener/complete"

	"github.com/grafana/tanka/pkg/tanka"
)

func fmtCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "fmt <FILES|DIRECTORIES>",
		Short: "format Jsonnet code",
		Args: cli.Args{
			Validator: cli.ValidateFunc(func(args []string) error {
				if len(args) == 0 {
					return errors.New("At least one file or directory is required")
				}
				return nil
			}),
			Predictor: complete.PredictFiles("*.*sonnet"),
		},
	}

	inplace := cmd.Flags().BoolP("inplace", "i", true, "save changes back to the original file instead of stdout")
	test := cmd.Flags().BoolP("test", "t", false, "exit with non-zero when changes would be made")
	exclude := cmd.Flags().StringSliceP("exclude", "e", []string{"**/.*", ".*", "**/vendor/**", "vendor/**"}, "globs to exclude")
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

		var outFn tanka.OutFn = nil
		switch {
		case *test:
			outFn = func(name, content string) error { return nil }
		case !*inplace:
			outFn = func(name, content string) error {
				fmt.Printf("// %s\n%s\n", name, content)
				return nil
			}
		}

		changed, err := tanka.Format(args, &tanka.FormatOpts{
			Excludes:   globs,
			OutFn:      outFn,
			PrintNames: *verbose,
		})
		if err != nil {
			return err
		}

		if *verbose {
			log.Println()
		}

		switch {
		case *test && len(changed) > 0:
			log.Println("The following files are not properly formatted:")
			for _, s := range changed {
				log.Println(s)
			}
			os.Exit(16)
		case len(changed) == 0:
			log.Println("All discovered files are already formatted. No changes were made")
		case len(changed) > 0:
			log.Printf("Reformatted %v files", len(changed))
		}

		return nil
	}

	return cmd
}
