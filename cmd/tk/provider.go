package main

import "github.com/spf13/cobra"

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "[Requires Provider] apply the configuration to the target",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {}
	return cmd
}

func diffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "[Requires Provider] print differences between the configuration and the target",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {}
	return cmd
}
