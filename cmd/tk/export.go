package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/spf13/cobra"

	"github.com/grafana/tanka/pkg/tanka"
)

func exportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <environment> <outputDir>",
		Short: "write each resources as a YAML file",
		Args:  cobra.ExactArgs(2),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	vars := workflowFlags(cmd.Flags())
	getExtCode := extCodeParser(cmd.Flags())
	format := cmd.Flags().String("format", "{{.apiVersion}}.{{.kind}}-{{.metadata.name}}", "https://tanka.dev/exporting#filenames")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// dir must be empty
		to := args[1]
		empty, err := dirEmpty(to)
		if err != nil {
			log.Fatalln("Checking target dir:", err)
		}
		if !empty {
			log.Fatalln("Target dir", to, "not empty. Aborting.")
		}

		// exit early if the template is bad
		tmpl, err := template.New("").Parse(*format)
		if err != nil {
			log.Fatalln("Parsing name format:", err)
		}

		// get the manifests
		res, err := tanka.Show(args[0],
			tanka.WithExtCode(getExtCode()),
			tanka.WithTargets(stringsToRegexps(vars.targets)...),
		)
		if err != nil {
			log.Fatalln(err)
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
				log.Fatalln("Writing manifest:", err)
			}
		}
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
