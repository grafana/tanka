package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-clix/cli"
	"github.com/posener/complete"

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
		Short: "export JSONNET_PATH for use with other jsonnet tools",
		Use:   "jpath [<file/dir>]",
		Args: cli.Args{
			Validator: cli.ValidateFunc(func(args []string) error {
				if len(args) != 1 {
					return errors.New("One file or directory is required")
				}
				return nil
			}),
			Predictor: complete.PredictFiles("*.*sonnet"),
		},
	}

	debug := cmd.Flags().BoolP("debug", "d", false, "show debug info")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		path := args[0]

		entrypoint, err := jpath.Entrypoint(path)
		if err != nil {
			return fmt.Errorf("Resolving JPATH: %s", err)
		}

		jsonnet_path, base, root, err := jpath.Resolve(entrypoint)
		if err != nil {
			return fmt.Errorf("Resolving JPATH: %s", err)
		}

		if *debug {
			// log to debug info to stderr
			log.Println("main:", entrypoint)
			log.Println("rootDir:", root)
			log.Println("baseDir:", base)
			log.Println("jpath:", jsonnet_path)
		}

		// print export JSONNET_PATH to stdout
		fmt.Printf("%s", strings.Join(jsonnet_path, ":"))

		return nil
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
