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

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

func toolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "handy utilities for working with jsonnet",
		Use:   "tool [command]",
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
		Use:   "imports <directory>",
		Short: "list all transitive imports of an environment",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
		Run: func(cmd *cobra.Command, args []string) {
			var modFiles []string
			if cmd.Flag("check").Changed {
				var err error
				modFiles, err = gitChangedFiles(cmd.Flag("check").Value.String())
				if err != nil {
					log.Fatalln("invoking git:", err)
				}
			}

			dir, err := filepath.Abs(args[0])
			if err != nil {
				log.Fatalln("Loading environment:", err)
			}

			fi, err := os.Stat(dir)
			if err != nil {
				log.Fatalln("Loading environment:", err)
			}

			if !fi.IsDir() {
				log.Fatalln("The argument must be an environment's directory, but this does not seem to be the case.")
			}

			deps, err := jsonnet.TransitiveImports(dir)
			if err != nil {
				log.Fatalln("Resolving imports:", err)
			}

			root, err := gitRoot()
			if err != nil {
				log.Fatalln("Invoking git:", err)
			}
			if modFiles != nil {
				for _, m := range modFiles {
					mod := filepath.Join(root, m)
					if err != nil {
						log.Fatalln(err)
					}

					for _, dep := range deps {
						if mod == dep {
							fmt.Printf("Rebuild required. File `%s` imports `%s`, which has been changed in `%s`.\n", args[0], dep, cmd.Flag("check").Value.String())
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

func gitRoot() (string, error) {
	s, err := git("rev-parse", "--show-toplevel")
	return strings.TrimRight(s, "\n"), err
}

func gitChangedFiles(sha string) ([]string, error) {
	f, err := git("diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return nil, err
	}
	return strings.Split(f, "\n"), nil
}

func git(argv ...string) (string, error) {
	cmd := exec.Command("git", argv...)
	cmd.Stderr = os.Stderr
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
