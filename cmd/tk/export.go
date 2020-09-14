package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/Masterminds/sprig/v3"
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

	format := cmd.Flags().String("format", "{{.apiVersion}}.{{.kind}}-{{.metadata.name}}", "https://tanka.dev/exporting#filenames")
	extension := cmd.Flags().String("extension", "yaml", "File extension")
	merge := cmd.Flags().Bool("merge", false, "Allow merging with existing directory")

	vars := workflowFlags(cmd.Flags())
	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		// dir must be empty
		to := args[1]
		empty, err := dirEmpty(to)
		if err != nil {
			return fmt.Errorf("Checking target dir: %s", err)
		}
		if !empty && !*merge {
			return fmt.Errorf("Output dir `%s` not empty. Pass --merge to ignore this", to)
		}

		// exit early if the template is bad

		// Replace all os.path separators in string with BelRune for creating subfolders
		replacedFormat := strings.Replace(*format, string(os.PathSeparator), BelRune, -1)

		tmpl, err := template.New("").
			Funcs(sprig.TxtFuncMap()). // register Masterminds/sprig
			Parse(replacedFormat)      // parse template
		if err != nil {
			return fmt.Errorf("Parsing name format: %s", err)
		}

		// get the manifests
		res, err := tanka.Show(args[0], tanka.Opts{
			JsonnetOpts: getJsonnetOpts(),
			Filters:     stringsToRegexps(vars.targets),
		})
		if err != nil {
			return err
		}

		// write each to a file
		for _, m := range res {
			buf := bytes.Buffer{}
			if err := tmpl.Execute(&buf, m); err != nil {
				log.Fatalln("executing name template:", err)
			}

			// Replace all os.path separators in string in order to not accidentally create subfolders
			name := strings.Replace(buf.String(), string(os.PathSeparator), "-", -1)
			// Replace the BEL character inserted with a path separator again in order to create a subfolder
			name = strings.Replace(name, BelRune, string(os.PathSeparator), -1)

			// Create all subfolders in path
			path := filepath.Join(to, name+"."+*extension)

			// Abort if already exists
			if exists, err := fileExists(path); err != nil {
				return err
			} else if exists {
				return fmt.Errorf("File '%s' already exists. Aborting", path)
			}

			// Write file
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				return fmt.Errorf("creating filepath '%s': %s", filepath.Dir(path), err)
			}
			data := m.String()
			if err := ioutil.WriteFile(path, []byte(data), 0644); err != nil {
				return fmt.Errorf("writing manifest: %s", err)
			}
		}

		return nil
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
