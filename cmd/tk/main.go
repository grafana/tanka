package main

import (
	"log"
	"os"

	"github.com/sh0rez/tanka/pkg/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sh0rez/tanka/pkg/config/v1alpha1"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
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
		fmtCmd(),
		debugCmd(),
	)

	// other commands
	rootCmd.AddCommand(completionCommand(rootCmd))

	// Configuration
	viper.SetConfigName("spec")
	viper.AddConfigPath(".")

	// handle deprecated ksonnet spec
	for old, new := range deprecated {
		viper.RegisterAlias(new, old)
	}

	// Configuration
	if err := viper.ReadInConfig(); err != nil {
		// just run fine without config. Apply and Diff won't work
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := rootCmd.Execute(); err != nil {
				log.Fatalln("Ouch:", err)
			}
			os.Exit(1)
		}

		log.Fatalln("Reading config:", err)
	}
	checkDeprecated()

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln("Parsing config:", err)
	}

	// Kubernetes
	kube = &config.Spec
	if err := kube.Init(); err != nil {
		log.Fatalln("initializing:", err)
	}

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
