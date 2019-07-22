package main

import (
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sh0rez/tanka/pkg/config/v1alpha1"
	"github.com/sh0rez/tanka/pkg/provider"
	"github.com/sh0rez/tanka/pkg/provider/kubernetes"
)

// Version is the current version of tk. To be set during build
var Version = "dev"

var (
	config    = &v1alpha1.Config{}
	prov      provider.Provider
	provName  string
	providers = map[string]provider.EmptyConstructor{
		"kubernetes": func() provider.Provider { return &kubernetes.Kubernetes{} },
	}
)

// list of deprecated config keys and their alternatives
// however, they still work and are aliased internally
var deprecated = map[string]string{
	"namespace": "spec.kubernetes.namespace",
	"server":    "spec.kubernetes.server",
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

	// provider commands
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
		// just run fine without config. Provider features won't work (apply, show, diff)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := rootCmd.Execute(); err != nil {
				log.Fatalln("Ouch:", err)
			}
			os.Exit(1)
		}

		log.Fatalln(err)
	}
	checkDeprecated()

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	// Provider
	var err error
	prov, provName, err = setupProvider()
	if err != nil {
		log.Fatalln(err)
	}

	rootCmd.AddCommand(providerCmd())

	// Run!
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(rootCmd.Execute())
	}
}

func setupProvider() (provider.Provider, string, error) {
	for name, construct := range providers {
		if cfg, ok := config.Spec[name]; ok {
			pro := construct()
			if err := mapstructure.Decode(cfg, &pro); err != nil {
				return nil, "", err
			}
			return pro, name, nil
		}
	}

	return nil, "none", nil
}

func checkDeprecated() {
	for old, use := range deprecated {
		if viper.IsSet(old) {
			log.Printf("Warning: `%s` is deprecated, use `%s` instead.", old, use)
		}
	}
}
