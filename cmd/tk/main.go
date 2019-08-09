package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/posener/complete"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/grafana/tanka/pkg/cmp"
	"github.com/grafana/tanka/pkg/config/v1alpha1"
	"github.com/grafana/tanka/pkg/kubernetes"
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
		log.Fatalln("Ouch:", rootCmd.Execute())
	}
}

func setupConfiguration(baseDir string) *v1alpha1.Config {
	viper.SetConfigName("spec")

	// if the baseDir arg is not a dir, abort
	pwd, err := filepath.Abs(baseDir)
	if err != nil {
		return nil
	}
	viper.AddConfigPath(pwd)

	// handle deprecated ksonnet spec
	for old, new := range deprecated {
		viper.RegisterAlias(new, old)
	}

	// read it
	if err := viper.ReadInConfig(); err != nil {
		// just run fine without config. Provider features won't work (apply, show, diff)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		log.Fatalln(err)
	}
	checkDeprecated()

	var config v1alpha1.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}
	return &config
}

func checkDeprecated() {
	for old, use := range deprecated {
		if viper.IsSet(old) {
			log.Printf("Warning: `%s` is deprecated, use `%s` instead.", old, use)
		}
	}
}
