package main

import (
	"encoding/json"
	"log"

	"github.com/spf13/cobra"

	"github.com/grafana/tanka/pkg/tanka"
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
		raw, _, err := tanka.Eval(args[0], nil)
		if err != nil {
			log.Fatalln(err, nil)
		}

		out, err := json.MarshalIndent(raw, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		pageln(string(out))
	}

	return cmd
}

// func eval(workdir string) (string, error) {
// 	// pwd, err := filepath.Abs(workdir)
// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// _, baseDir, _, err := jpath.Resolve(pwd)
// 	// if err != nil {
// 	// 	return "", errors.Wrap(err, "resolving jpath")
// 	// }
// 	// json, err := jsonnet.EvaluateFile(filepath.Join(baseDir, "main.jsonnet"))
// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	return json, nil
// }
