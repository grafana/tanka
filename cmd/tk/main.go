package main

import (
	"fmt"
	"os"
	"runtime/pprof"
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

	// workflow commands
	addCommandsWithLogLevelOption(
		rootCmd,
		applyCmd(),
		showCmd(),
		diffCmd(),
		pruneCmd(),
		deleteCmd(),
	)

	addCommandsWithLogLevelOption(
		rootCmd,
		envCmd(),
		statusCmd(),
		exportCmd(),
	)

	// jsonnet commands
	addCommandsWithLogLevelOption(
		rootCmd,
		fmtCmd(),
		lintCmd(),
		evalCmd(),
		initCmd(),
		toolCmd(),
	)

	// external commands prefixed with "tk-"
	addCommandsWithLogLevelOption(
		rootCmd,
		prefixCommands("tk-")...,
	)

	outputPprofFile := os.Getenv("TK_PPROF_FILE")

	// So that the other defers can run
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	if outputPprofFile != "" {
		pprofFile, err := os.Create(outputPprofFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating pprof file: %s\n", err)
			exitCode = 1
			return
		}
		defer pprofFile.Close()

		err = pprof.StartCPUProfile(pprofFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error starting pprof: %s\n", err)
			exitCode = 1
			return
		}
		defer pprof.StopCPUProfile()
	}

	// Run!
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		exitCode = 1
		return
	}
}

func addCommandsWithLogLevelOption(rootCmd *cli.Command, cmds ...*cli.Command) {
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
