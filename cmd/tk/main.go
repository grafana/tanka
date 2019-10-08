package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/posener/complete"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/grafana/tanka/pkg/cmp"
	"github.com/grafana/tanka/pkg/kubernetes"
	"github.com/grafana/tanka/pkg/spec"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// Version is the current version of the tk command.
// To be overwritten at build time
var Version = "dev"

// primary handlers
var (
	config = &v1alpha1.Config{}
	kube   *kubernetes.Kubernetes
)

// describing variables
var (
	verbose     = false
	interactive = terminal.IsTerminal(int(os.Stdout.Fd()))
)

// list of deprecated config keys and their alternatives
// however, they still work and are aliased internally
var deprecated = map[string]string{
	"namespace": "spec.namespace",
	"server":    "spec.apiServer",
	"team":      "metadata.labels.team",
}

func main() {
	log.SetFlags(0)
	rootCmd := &cobra.Command{
		Use:              "tk",
		Short:            "tanka <3 jsonnet",
		Version:          Version,
		TraverseChildren: true,
		// Configuration
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				return
			}
			config = setupConfiguration(args[0])
			if config == nil {
				return
			}

			// Kubernetes
			kube = kubernetes.New(config.Spec)

		},
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
	config, err := spec.ParseDir(baseDir)
	if err != nil {
		switch err.(type) {
		// just run fine without config. Provider features won't work (apply, show, diff)
		case viper.ConfigFileNotFoundError:
			return nil
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
