package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

func toolCmd() *cli.Command {
	cmd := &cli.Command{
		Short: "handy utilities for working with jsonnet",
		Use:   "tool [command]",
	}
	cmd.AddCommand(
		jpathCmd(),
		importsCmd(),
		chartsCmd(),
	)
	return cmd
}

func jpathCmd() *cli.Command {
	cmd := &cli.Command{
		Short: "print information about the jpath",
		Use:   "jpath",
		Run: func(cmd *cli.Command, args []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			path, base, root, err := jpath.Resolve(pwd)
			if err != nil {
				return fmt.Errorf("Resolving JPATH: %s", err)
			}
			entrypoint, err := jpath.GetEntrypoint(base)
			if err != nil {
				return fmt.Errorf("Resolving JPATH: %s", err)
			}
			fmt.Println("main:", entrypoint)
			fmt.Println("rootDir:", root)
			fmt.Println("baseDir:", base)
			fmt.Println("jpath:", path)

			return nil
		},
	}
	return cmd
}

func importsCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "imports <directory>",
		Short: "list all transitive imports of an environment",
		Args:  workflowArgs,
	}

	check := cmd.Flags().StringP("check", "c", "", "git commit hash to check against")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		var modFiles []string
		if *check != "" {
			var err error
			modFiles, err = gitChangedFiles(*check)
			if err != nil {
				return fmt.Errorf("invoking git: %s", err)
			}
		}

		dir, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("Loading environment: %s", err)
		}

		fi, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("Loading environment: %s", err)
		}

		if !fi.IsDir() {
			return fmt.Errorf("The argument must be an environment's directory, but this does not seem to be the case.")
		}

		deps, err := jsonnet.TransitiveImports(dir)
		if err != nil {
			return fmt.Errorf("Resolving imports: %s", err)
		}

		root, err := gitRoot()
		if err != nil {
			return fmt.Errorf("Invoking git: %s", err)
		}
		if modFiles != nil {
			for _, m := range modFiles {
				mod := filepath.Join(root, m)
				if err != nil {
					return err
				}

				for _, dep := range deps {
					if mod == dep {
						fmt.Printf("Rebuild required. File `%s` imports `%s`, which has been changed in `%s`.\n", args[0], dep, *check)
						os.Exit(16)
					}
				}
			}
			fmt.Printf("Rebuild not required, because no imported files have been changed in `%s`.\n", *check)
			os.Exit(0)
		}

		s, err := json.Marshal(deps)
		if err != nil {
			return fmt.Errorf("Formatting: %s", err)
		}
		fmt.Println(string(s))

		return nil
	}

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
