package main

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/go-clix/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/grafana/tanka/pkg/tanka"
)

// describing variables
var (
	verbose     = false
	interactive = terminal.IsTerminal(int(os.Stdout.Fd()))
)

func main() {
	log.SetFlags(0)

	rootCmd := &cli.Command{
		Use:     "tk",
		Short:   "tanka <3 jsonnet",
		Version: tanka.CURRENT_VERSION,
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
