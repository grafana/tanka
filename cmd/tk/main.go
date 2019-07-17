package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/spf13/cobra"

	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/sh0rez/tanka/pkg/native"
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

	cmd.Run = func(cmd *cobra.Command, args []string) {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		jPath, baseDir, _ := jpath.Resolve(pwd)
		importer := jsonnet.FileImporter{
			JPaths: jPath,
		}

		vm := jsonnet.MakeVM()
		vm.Importer(&importer)
		for _, f := range native.Funcs() {
			vm.NativeFunction(f)
		}

		jsonnetFile := filepath.Join(baseDir, "main.jsonnet")
		jsonnetBytes, err := ioutil.ReadFile(jsonnetFile)
		if err != nil {
			panic(err)
		}

		jsonStr, err := vm.EvaluateSnippet(jsonnetFile, string(jsonnetBytes))
		if err != nil {
			panic(err)
		}

		fmt.Print(jsonStr)
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
