package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/posener/complete"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/grafana/tanka/pkg/cli/cmp"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Version is the current version of the tk command.
// To be overwritten at build time
var Version = "dev"

// describing variables
var (
	verbose     = false
	interactive = terminal.IsTerminal(int(os.Stdout.Fd()))
)

func main() {
	log.SetFlags(0)
	rootCmd := &cobra.Command{
		Use:              "tk",
		Short:            "tanka <3 jsonnet",
		Version:          Version,
		TraverseChildren: true,
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "")

	// Subcommands
	cobra.EnableCommandSorting = false

	// workflow commands
	rootCmd.AddCommand(
		applyCmd(),
		showCmd(),
		diffCmd(),
	)

	rootCmd.AddCommand(
		envCmd(),
		statusCmd(),
	)

	// jsonnet commands
	rootCmd.AddCommand(
		evalCmd(),
		initCmd(),
		toolCmd(),
	)

	// completion
	cmp.Handlers.Add("baseDir", complete.PredictFunc(
		func(complete.Args) []string {
			return findBaseDirs()
		},
	))

	c := complete.New("tk", cmp.Create(rootCmd))
	c.InstallName = "install-completion"
	c.UninstallName = "uninstall-completion"
	fs := &flag.FlagSet{}
	c.AddFlags(fs)
	rootCmd.Flags().AddGoFlagSet(fs)

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if c.Complete() {
			return
		}
		_ = cmd.Help()
	}

	// Run!
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Ouch:", err)
	}
}

func setupConfiguration(baseDir string) *v1alpha1.Config {
	_, baseDir, rootDir, err := jpath.Resolve(baseDir)
	if err != nil {
		log.Fatalln("Resolving jpath:", err)
	}

	// name of the environment: relative path from rootDir
	name, _ := filepath.Rel(rootDir, baseDir)

	config, err := spec.ParseDir(baseDir, name)
	if err != nil {
		switch err.(type) {
		// the config includes deprecated fields
		case spec.ErrDeprecated:
			if verbose {
				fmt.Print(err)
			}
		// some other error
		default:
			log.Fatalf("Reading spec.json: %s", err)
		}
	}

	return config
}
