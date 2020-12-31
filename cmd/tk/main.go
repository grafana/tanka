package main

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/go-clix/cli"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
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

	// Run!
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(color.RedString("Error:"), err)
	}
}

func setupConfiguration(baseDir string) *v1alpha1.Environment {
	env, err := tanka.Load(baseDir, tanka.Opts{
		JsonnetOpts: tanka.JsonnetOpts{EvalScript: tanka.EnvsOnlyEvalScript},
	})
	if err != nil {
		log.Fatalln(err)
	}

	return env.Env
}
