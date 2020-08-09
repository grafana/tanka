package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/pflag"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/tanka"
)

func evalCmd() *cli.Command {
	cmd := &cli.Command{
		Short: "evaluate the jsonnet to json",
		Use:   "eval <path>",
		Args:  workflowArgs,
	}

	getExtCode := extCodeParser(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		raw, err := tanka.Eval(args[0],
			tanka.WithExtCode(getExtCode()),
			tanka.WithMainfile("main.jsonnet"),
		)

		if err != nil {
			return err
		}

		out, err := json.MarshalIndent(raw, "", "  ")
		if err != nil {
			return err
		}

		pageln(string(out))
		return nil
	}

	return cmd
}

func extCodeParser(fs *pflag.FlagSet) func() map[string]string {
	// need to use StringArray instead of StringSlice, because pflag attempts to
	// parse StringSlice using the csv parser, which breaks when passing objects
	values := fs.StringArrayP("extCode", "e", nil, "Inject any Jsonnet from the outside (Format: key=<code>)")
	strs := fs.StringArray("extVar", nil, "Inject a string from the outside (Format: key=value)")

	return func() map[string]string {
		m := make(map[string]string)
		for _, s := range *values {
			split := strings.SplitN(s, "=", 2)
			if len(split) != 2 {
				log.Fatalf("extCode argument has wrong format: `%s`. Expected `key=<code>`", s)
			}
			m[split[0]] = split[1]
		}

		for _, s := range *strs {
			split := strings.SplitN(s, "=", 2)
			if len(split) != 2 {
				log.Fatalf("extCode argument has wrong format: `%s`. Expected `key=<value>`", s)
			}
			m[split[0]] = fmt.Sprintf(`"%s"`, split[1])
		}
		return m
	}
}
