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

	"github.com/grafana/tanka/pkg/cli"
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

	vars := workflowFlags(cmd.Flags())
	getExtCode := extCodeParser(cmd.Flags())
	format := cmd.Flags().String("format", "{{.apiVersion}}.{{.kind}}-{{.metadata.name}}", "https://tanka.dev/exporting#filenames")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		// dir must be empty
		to := args[1]
		empty, err := dirEmpty(to)
		if err != nil {
			return fmt.Errorf("Checking target dir: %s", err)
		}
		if !empty {
			return fmt.Errorf("Target dir `%s` not empty. Aborting.", to)
		}

		// exit early if the template is bad
		tmpl, err := template.New("").Parse(*format)
		if err != nil {
			return fmt.Errorf("Parsing name format: %s", err)
		}

		// get the manifests
		res, err := tanka.Show(args[0],
			tanka.WithExtCode(getExtCode()),
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
		)
		if err != nil {
			return err
		}

		// write each to a file
		for _, m := range res {
			buf := bytes.Buffer{}
			if err := tmpl.Execute(&buf, m); err != nil {
				log.Fatalln("executing name template:", err)
			}
			name := strings.Replace(buf.String(), "/", "-", -1)

			data := m.String()
			if err := ioutil.WriteFile(filepath.Join(to, name+".yaml"), []byte(data), 0644); err != nil {
				return fmt.Errorf("Writing manifest: %s", err)
			}
		}

		return nil
	}
	return cmd
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
