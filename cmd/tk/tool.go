package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/sh0rez/tanka/pkg/jsonnet"
)

func toolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "handy utilities for working with jsonnet",
		Use:   "tool",
	}
	cmd.AddCommand(jpathCmd())
	cmd.AddCommand(importsCmd())
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
			path, base, root, err := jpath.Resolve(pwd)
			if err != nil {
				log.Fatalln("Resolving JPATH:", err)
			}
			fmt.Println("main:", filepath.Join(base, "main.jsonnet"))
			fmt.Println("rootDir:", root)
			fmt.Println("baseDir:", base)
			fmt.Println("jpath:", path)
			return nil
		},
	}
	return cmd
}

func importsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "imports [file]",
		Short: "list all transitive imports of a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var modFiles []string
			if cmd.Flag("check").Changed {
				argv := []string{"diff-tree", "--no-commit-id", "--name-only", "-r", cmd.Flag("check").Value.String()}
				c := exec.Command("git", argv...)
				c.Stderr = os.Stderr
				var buf bytes.Buffer
				c.Stdout = &buf
				if err := c.Run(); err != nil {
					log.Fatalln("Invoking git:", err)
				}
				modFiles = strings.Split(buf.String(), "\n")
			}

			f, err := filepath.Abs(args[0])
			if err != nil {
				log.Fatalln("Opening file:", err)
			}

			deps, err := jsonnet.TransitiveImports(f)
			if err != nil {
				log.Fatalln("resolving imports:", err)
			}

			// include main.jsonnet as well
			deps = append(deps, f)

			if modFiles != nil {
				for _, m := range modFiles {
					for _, d := range deps {
						if m == d {
							fmt.Printf("Rebuild required. File `%s` imports `%s`, which has been changed in `%s`.\n", args[0], d, cmd.Flag("check").Value.String())
							os.Exit(16)
						}
					}
				}
				fmt.Printf("Rebuild not required, because no imported files have been changed in `%s`.\n", cmd.Flag("check").Value.String())
				os.Exit(0)
			}

			s, err := json.Marshal(deps)
			if err != nil {
				log.Fatalln("Formatting:", err)
			}
			fmt.Println(string(s))
		},
	}

	cmd.Flags().StringP("check", "c", "", "git commit hash to check against")

	return cmd
}
