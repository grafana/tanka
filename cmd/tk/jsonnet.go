package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/tanka"
)

func evalCmd() *cli.Command {
	cmd := &cli.Command{
		Short: "evaluate the jsonnet to json",
		Use:   "eval <path>",
		Args:  workflowArgs,
	}

	evalPattern := cmd.Flags().StringP("eval", "e", "", "Evaluate expression on output of jsonnet")

	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		jsonnetOpts := tanka.Opts{
			JsonnetOpts: getJsonnetOpts(),
		}
		if *evalPattern != "" {
			jsonnetOpts.EvalScript = fmt.Sprintf(tanka.PatternEvalScript, *evalPattern)
		}
		raw, err := tanka.Eval(args[0], jsonnetOpts)

		if raw == nil && err != nil {
			return err
		}

		out, err := json.MarshalIndent(raw, "", "  ")
		if err != nil {
			return err
		}

		if err := pageln(string(out)); err != nil {
			return err
		}

		return nil
	}

	return cmd
}
