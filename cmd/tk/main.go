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
	outputPprofFile := os.Getenv("TANKA_PPROF_FILE")

	// So that the other defers can run
	exitCode, exitMessage := 0, ""
	defer func() {
		if exitMessage != "" {
			fmt.Fprintln(os.Stderr, exitMessage)
		}
		os.Exit(exitCode)
	}()
	exitF := func(code int, messageF string, args ...interface{}) {
		exitCode = code
		exitMessage = fmt.Sprintf(messageF, args...)
	}

	if outputPprofFile != "" {
		pprofFile, err := os.Create(outputPprofFile)
		if err != nil {
			exitF(2, "Error creating pprof file: %s\n", err)
			return
		}
		defer pprofFile.Close()

		err = pprof.StartCPUProfile(pprofFile)
		if err != nil {
			exitF(2, "Error starting pprof: %s\n", err)
			return
		}
		defer pprof.StopCPUProfile()
	}

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

	// Run!
	if err := rootCmd.Execute(); err != nil {
		exitF(1, "Error: %s\n", err)
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
