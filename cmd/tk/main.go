package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"github.com/rs/zerolog"
	"golang.org/x/term"

	"github.com/grafana/tanka/pkg/tanka"
)

var interactive = term.IsTerminal(int(os.Stdout.Fd()))

func main() {
	rootCmd := &cli.Command{
		Use:     "tk",
		Short:   "tanka <3 jsonnet",
		Version: tanka.CurrentVersion,
	}

	addCommandsWithLogLevel := func(cmds ...*cli.Command) {
		for _, cmd := range cmds {
			levels := []string{zerolog.Disabled.String(), zerolog.FatalLevel.String(), zerolog.ErrorLevel.String(), zerolog.WarnLevel.String(), zerolog.InfoLevel.String(), zerolog.DebugLevel.String(), zerolog.TraceLevel.String()}
			cmd.Flags().String("log-level", zerolog.InfoLevel.String(), "possible values: "+strings.Join(levels, ", "))

			cmdRun := cmd.Run
			cmd.Run = func(cmd *cli.Command, args []string) error {
				level, err := zerolog.ParseLevel(cmd.Flags().Lookup("log-level").Value.String())
				if err != nil {
					return err
				}
				zerolog.SetGlobalLevel(level)

				return cmdRun(cmd, args)
			}
			rootCmd.AddCommand(cmd)
		}
	}

	// workflow commands
	addCommandsWithLogLevel(
		applyCmd(),
		showCmd(),
		diffCmd(),
		pruneCmd(),
		deleteCmd(),
	)

	addCommandsWithLogLevel(
		envCmd(),
		statusCmd(),
		exportCmd(),
	)

	// jsonnet commands
	addCommandsWithLogLevel(
		fmtCmd(),
		lintCmd(),
		evalCmd(),
		initCmd(),
		toolCmd(),
	)

	// external commands prefixed with "tk-"
	addCommandsWithLogLevel(
		prefixCommands("tk-")...,
	)

	// Run!
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
