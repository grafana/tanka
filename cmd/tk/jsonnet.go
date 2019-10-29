package main

import (
	"log"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func evalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "evaluate the jsonnet to json",
		Use:   "eval <path>",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		json, err := eval(args[0])
		if err != nil {
			log.Fatalln("evaluating:", err)
		}
		pageln(json)
	}

	return cmd
}

func eval(workdir string) (string, error) {
	pwd, err := filepath.Abs(workdir)
	if err != nil {
		return "", err
	}
	_, baseDir, _, err := jpath.Resolve(pwd)
	if err != nil {
		return "", errors.Wrap(err, "resolving jpath")
	}
	json, err := jsonnet.EvaluateFile(filepath.Join(baseDir, "main.jsonnet"))
	if err != nil {
		return "", err
	}
	return json, nil
}
