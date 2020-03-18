package main

import (
	"fmt"
	"os"

	"github.com/go-clix/cli"
	"github.com/gobwas/glob"

	"github.com/grafana/tanka/pkg/tanka"
)

func fmtCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "fmt <FILES|DIRECTORIES>",
		Short: "format Jsonnet code",
		Args:  cli.ArgsAny(),
	}

	inplace := cmd.Flags().BoolP("inplace", "i", true, "save changes back to the original file instead of stdout")
	test := cmd.Flags().BoolP("test", "t", false, "exit with non-zero when changes would be made")
	exclude := cmd.Flags().StringSliceP("exclude", "e", []string{"**/.*", ".*", "**/vendor/**", "vendor/**"}, "globs to exclude")

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

		err := tanka.Format(args, &tanka.FormatOpts{
			Excludes: globs,
			Test:     *test,
			OutFn:    outFn,
		})
		if _, ok := err.(tanka.ErrorNotFormatted); ok {
			fmt.Println(err)
			os.Exit(16)
		} else if err != nil {
			fmt.Printf("%T\n", err)
			return err
		}

		return nil
	}

	return cmd
}
