package main

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sh0rez/tanka/pkg/config/v1alpha1"
	"github.com/sh0rez/tanka/pkg/kubernetes"
)

// Version is the current version of the tk command.
// To be overwritten at build time
var Version = "dev"

var (
	config = &v1alpha1.Config{}
	kube   *kubernetes.Kubernetes
)

// list of deprecated config keys and their alternatives
// however, they still work and are aliased internally
var deprecated = map[string]string{
	"namespace": "spec.namespace",
	"server":    "spec.apiServer",
	"team":      "metadata.labels.team",
}

func main() {
	rootCmd := &cobra.Command{
		Use:              "tk",
		Short:            "tanka <3 jsonnet",
		Version:          Version,
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Configuration parsing. Note this is using return to abort actions
			viper.SetConfigName("spec")

			// no args = no command that has a baseDir passed, abort
			if len(args) == 0 {
				return
			}

			// if the first arg is not a dir, abort
			pwd, err := filepath.Abs(args[0])
			if err != nil {
				return
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
					return
				}

				log.Fatalln(err)
			}
			checkDeprecated()

			if err := viper.Unmarshal(&config); err != nil {
				log.Fatalln(err)
			}

			// Kubernetes
			kube = kubernetes.New(config.Spec)

		},
	}
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "")

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
		debugCmd(),
	)

	// other commands
	rootCmd.AddCommand(completionCommand(rootCmd))

	// Run!
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Ouch:", rootCmd.Execute())
	}
}
func checkDeprecated() {
	for old, use := range deprecated {
		if viper.IsSet(old) {
			log.Printf("Warning: `%s` is deprecated, use `%s` instead.", old, use)
		}
	}
}
