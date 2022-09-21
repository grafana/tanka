package main

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/go-clix/cli"
	"golang.org/x/term"

	"github.com/grafana/tanka/pkg/tanka"
)

var interactive = term.IsTerminal(int(os.Stdout.Fd()))

func main() {
	log.SetFlags(0)

	rootCmd := &cli.Command{
		Use:     "tk",
		Short:   "tanka <3 jsonnet",
		Version: tanka.CurrentVersion,
	}

	// workflow commands
	rootCmd.AddCommand(
		applyCmd(),
		showCmd(),
		diffCmd(),
		pruneCmd(),
		deleteCmd(),
	)

	rootCmd.AddCommand(
		envCmd(),
		statusCmd(),
		exportCmd(),
	)

	// jsonnet commands
	rootCmd.AddCommand(
		fmtCmd(),
		lintCmd(),
		evalCmd(),
		initCmd(),
		toolCmd(),
	)

	// external commands prefixed with "tk-"
	rootCmd.AddCommand(
		prefixCommands("tk-")...,
	)

	// Run!
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(color.RedString("Error:"), err)
	}
}
