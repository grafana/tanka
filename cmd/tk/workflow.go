package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply the configuration to the cluster",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict()
		if err != nil {
			log.Fatalln("evaluating jsonnet:", err)
		}

		desired, err := kube.Reconcile(raw)
		if err != nil {
			log.Fatalln("reconciling:", err)
		}

		if err := kube.Apply(desired); err != nil {
			log.Fatalln("applying:", err)
		}
	}
	return cmd
}

func diffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "differences between the configuration and the cluster",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict()
		if err != nil {
			log.Fatalln("evaluating jsonnet:", err)
		}

		desired, err := kube.Reconcile(raw)
		if err != nil {
			log.Fatalln("reconciling:", err)
		}

		changes, err := kube.Diff(desired)
		if err != nil {
			log.Fatalln("diffing:", err)
		}

		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			if err := quick.Highlight(os.Stdout, changes, "diff", "terminal", "vim"); err != nil {
				log.Fatalln("highlighting:", err)
			}
		} else {
			fmt.Println(changes)
		}
	}
	return cmd
}

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "jsonnet as yaml",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict()
		if err != nil {
			log.Fatalln("evaluating jsonnet:", err)
		}

		state, err := kube.Reconcile(raw)
		if err != nil {
			log.Fatalln("reconciling:", err)
		}

		pretty, err := kube.Fmt(state)
		if err != nil {
			log.Fatalln("pretty printing state:", err)
		}
		fmt.Println(pretty)
	}
	return cmd
}
