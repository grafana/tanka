package main

import (
	"io"
	"os"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/tanka"
)

// BelRune is a string of the Ascii character BEL which made computers ring in ancient times
// We use it as "magic" char for the subfolder creation as it is a non printable character and thereby will never be
// in a valid filepath by accident. Only when we include it.
const BelRune = string(rune(7))

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
		return tanka.ExportEnvironment(args[0], args[1], &opts)
	}
	return cmd
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func dirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		return true, os.MkdirAll(dir, os.ModePerm)
	} else if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
