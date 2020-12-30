package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/fatih/color"
	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/spec"
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
		// no spec.json is found, try parsing main.jsonnet
		case spec.ErrNoSpec:
			_, config, err := tanka.ParseEnv(baseDir, tanka.ParseOpts{Evaluator: tanka.EnvsOnlyEvaluator})
			if err != nil {
				switch err.(type) {
				case tanka.ErrNoEnv:
					return nil
				default:
					log.Fatalf("Reading main.jsonnet: %s", err)
				}
			}
			return config
		// some other error
		default:
			log.Fatalf("Reading spec.json: %s", err)
		}
	}

	return config
}
