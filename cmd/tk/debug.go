package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/spf13/cobra"
)

func debugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "debug utilities",
		Use:   "debug",
	}
	cmd.AddCommand(jpathCmd())
	return cmd
}

func jpathCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "print information about the jpath",
		Use:   "jpath [directory]",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
		Run: func(cmd *cobra.Command, args []string) {
			pwd, err := filepath.Abs(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			path, base, root, err := jpath.Resolve(pwd)
			if err != nil {
				log.Fatalln("resolving jpath:", err)
			}
			fmt.Println("main:", filepath.Join(base, "main.jsonnet"))
			fmt.Println("rootDir:", root)
			fmt.Println("baseDir:", base)
			fmt.Println("jpath:", path)
		},
	}
	return cmd
}
