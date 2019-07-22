package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func providerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "interact with providers",
	}
	cmd.AddCommand(providerListCmd())

	proCmd := prov.Cmd()
	proCmd.Use = provName
	cmd.AddCommand(proCmd)
	return cmd
}

func providerListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "print all available providers",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Available providers:", listProviders())
		},
	}
	return cmd
}

func listProviders() []string {
	keys := make([]string, len(providers))

	i := 0
	for k := range providers {
		keys[i] = k
		i++
	}
	return keys
}

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "[Requires Provider] apply the configuration to the target",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
	}
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

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "[Requires Provider] print the jsonnet in the target state format",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict()
		if err != nil {
			log.Fatalln("evaluating jsonnet:", err)
		}

		state, err := prov.Format(rawDict)
		if err != nil {
			log.Fatalln("invoking provider:", err)
		}

		pretty, err := prov.Show(state)
		if err != nil {
			log.Fatalln("pretty printing state:", err)
		}
		fmt.Println(pretty)
	}
	return cmd
}
