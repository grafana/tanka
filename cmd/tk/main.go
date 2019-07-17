package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/sh0rez/tanka/pkg/sonnet"
)

var Version = "dev"

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
		diffCmd(),
	)

	// jsonnet commands
	rootCmd.AddCommand(
		evalCmd(),
		fmtCmd(),
		debugCmd(),
	)

	// other commands
	rootCmd.AddCommand(completionCommand(rootCmd))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func fmtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "format .jsonnet and .libsonnet files",
		Use:   "fmt",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {}
	return cmd
}

func evalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "evaluate the jsonnet to json",
		Use:   "eval",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		_, baseDir, _ := jpath.Resolve(pwd)
		json, err := sonnet.EvaluateFile(filepath.Join(baseDir, "main.jsonnet"))
		if err != nil {
			return err
		}

		fmt.Print(json)
		return nil
	}

	return cmd
}

func completionCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use: "completion bash|zsh",
		Example: `
  eval "$(tk completion bash)"
  eval "$(tk completion zsh)"
		`,
		Short: `create bash/zsh auto-completion.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch shell := strings.ToLower(args[0]); shell {
			case "bash":
				return rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				return rootCmd.GenZshCompletion(os.Stdout)
			default:
				return fmt.Errorf("unknown shell %q. Only bash and zsh are supported", shell)
			}
		},
	}
}
