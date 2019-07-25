package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/spf13/cobra"
)

func completionCommand(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: "completion bash|zsh",
		Example: `
  eval "$(tk completion bash)"
  eval "$(tk completion zsh)"
		`,
		Short:     `create bash/zsh auto-completion.`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh"},
		Run: func(cmd *cobra.Command, args []string) {
			switch shell := strings.ToLower(args[0]); shell {
			case "bash":
				if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
					log.Fatalln(err)
				}
			case "zsh":
				if err := rootCmd.GenZshCompletion(os.Stdout); err != nil {
					log.Fatalln(err)
				}
			}
		},
	}
	cmd.AddCommand(baseDirsCommand())
	return cmd
}

func baseDirsCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "base-dirs",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			dirs := findBaseDirs()
			fmt.Println(strings.Join(dirs, " "))
		},
	}
}

// findBaseDirs searches for possible environments
func findBaseDirs() (dirs []string) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	_, _, _, err = jpath.Resolve(pwd)
	if err == jpath.ErrorNoRoot {
		return
	}

	if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(filepath.Join(path, "main.jsonnet")); err == nil {
			dirs = append(dirs, path)
		}
		return nil
	}); err != nil {
		log.Fatalln(err)
	}
	return dirs
}
