package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

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

	getExtCode := extCodeParser(cmd.Flags())

	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := tanka.Eval(args[0],
			tanka.WithExtCode(getExtCode()),
		)

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
				log.Fatalln("extCode argument has wrong format:", s+".", "Expected 'key=<code>'")
			}
			m[split[0]] = split[1]
		}

		for _, s := range *strs {
			split := strings.SplitN(s, "=", 2)
			if len(split) != 2 {
				log.Fatalln("extVar argument has wrong format:", s+".", "Expected 'key=value'")
			}
			m[split[0]] = fmt.Sprintf(`"%s"`, split[1])
		}
		return m
	}
}
