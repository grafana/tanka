package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
