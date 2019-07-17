package main

import (
	"fmt"
	"os"
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
		Use:   "jpath",
		RunE: func(cmd *cobra.Command, args []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			path, base, root := jpath.Resolve(pwd)
			fmt.Println("main:", filepath.Join(base, "main.jsonnet"))
			fmt.Println("rootDir:", root)
			fmt.Println("baseDir:", base)
			fmt.Println("jpath:", path)
			return nil
		},
	}
	return cmd
}
