package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/go-clix/cli"
	"github.com/pkg/errors"
	"github.com/posener/complete"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/grafana/tanka/pkg/tanka"
	"github.com/grafana/tanka/pkg/term"
)

func envCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "env [action]",
		Short: "manipulate environments",
	}

	cmd.AddCommand(
		envAddCmd(),
		envSetCmd(),
		envListCmd(),
		envRemoveCmd(),
	)

	return cmd
}

var kubectlContexts = cli.PredictFunc(
	func(complete.Args) []string {
		c, _ := client.Contexts()
		return c
	},
)

func envSetCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "set",
		Short: "update properties of an environment",
		Args:  workflowArgs,
		Predictors: complete.Flags{
			"server-from-context": kubectlContexts,
		},
	}

	// flags
	tmp := v1alpha1.Environment{}
	envSettingsFlags(&tmp, cmd.Flags())

	// removed name flag
	name := cmd.Flags().String("name", "", "")
	_ = cmd.Flags().MarkHidden("name")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if *name != "" {
			return fmt.Errorf("It looks like you attempted to rename the environment using `--name`. However, this is not possible with Tanka, because the environments name is inferred from the directories name. To rename the environment, rename its directory instead.")
		}

		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		if cmd.Flags().Changed("server-from-context") {
			server, err := client.IPFromContext(tmp.Spec.APIServer)
			if err != nil {
				return fmt.Errorf("Resolving IP from context: %s", err)
			}
			tmp.Spec.APIServer = server
		}

		cfg, err := tanka.Peek(path, tanka.JsonnetOpts{})
		if err != nil {
			return err
		}

		if tmp.Spec.APIServer != "" && tmp.Spec.APIServer != cfg.Spec.APIServer {
			fmt.Printf("updated spec.apiServer (`%s -> `%s`)\n", cfg.Spec.APIServer, tmp.Spec.APIServer)
			cfg.Spec.APIServer = tmp.Spec.APIServer
		}
		if tmp.Spec.Namespace != "" && tmp.Spec.Namespace != cfg.Spec.Namespace {
			fmt.Printf("updated spec.namespace (`%s -> `%s`)\n", cfg.Spec.Namespace, tmp.Spec.Namespace)
			cfg.Spec.Namespace = tmp.Spec.Namespace
		}
		if tmp.Spec.DiffStrategy != "" && tmp.Spec.DiffStrategy != cfg.Spec.DiffStrategy {
			fmt.Printf("updated spec.diffStrategy (`%s -> `%s`)\n", cfg.Spec.DiffStrategy, tmp.Spec.DiffStrategy)
			cfg.Spec.DiffStrategy = tmp.Spec.DiffStrategy
		}

		if err := writeJSON(cfg, filepath.Join(path, "spec.json")); err != nil {
			return err
		}

		return nil
	}
	return cmd
}

func envAddCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "add <path>",
		Short: "create a new environment",
		Args:  cli.ArgsExact(1),
	}
	cfg := v1alpha1.New()
	envSettingsFlags(cfg, cmd.Flags())
	inline := cmd.Flags().BoolP("inline", "i", false, "create an inline environment")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		if cmd.Flags().Changed("server-from-context") {
			server, err := client.IPFromContext(cfg.Spec.APIServer)
			if err != nil {
				return fmt.Errorf("Resolving IP from context: %s", err)
			}
			cfg.Spec.APIServer = server
		}

		if err := addEnv(args[0], cfg, *inline); err != nil {
			return err
		}

		return nil
	}
	return cmd
}

// used by initCmd() as well
func addEnv(dir string, cfg *v1alpha1.Environment, inline bool) error {
	path, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err != nil {
		// folder does not exist
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return errors.Wrap(err, "creating directory")
			}
		} else {
			// it exists
			if os.IsExist(err) {
				return fmt.Errorf("directory %s already exists", path)
			}
			// we have another error
			return errors.Wrap(err, "creating directory")
		}
	}

	rootDir, err := jpath.FindRoot(path)
	if err != nil {
		return err
	}
	// the other properties are already set by v1alpha1.New() and pflag.Parse()
	cfg.Metadata.Name, _ = filepath.Rel(rootDir, path)

	if inline {
		cfg.Data = struct{}{}
		// write main.jsonnet with inline tanka.dev/Environment
		if err := writeJsonnet(cfg, filepath.Join(path, "main.jsonnet")); err != nil {
			return err
		}
	} else {
		// write spec.json
		if err := writeJSON(cfg, filepath.Join(path, "spec.json")); err != nil {
			return err
		}

		// write main.jsonnet
		if err := writeJsonnet(struct{}{}, filepath.Join(path, "main.jsonnet")); err != nil {
			return err
		}
	}

	return nil
}

func envRemoveCmd() *cli.Command {
	return &cli.Command{
		Use:     "remove <path>",
		Aliases: []string{"rm"},
		Short:   "delete an environment",
		Args:    workflowArgs,
		Run: func(cmd *cli.Command, args []string) error {
			for _, arg := range args {
				path, err := filepath.Abs(arg)
				if err != nil {
					return fmt.Errorf("parsing environments name: %s", err)
				}
				if err := term.Confirm(fmt.Sprintf("Permanently removing the environment located at '%s'.", path), "yes"); err != nil {
					return err
				}
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("Removing '%s': %s", path, err)
				}
				fmt.Println("Removed", path)
			}
			return nil
		},
	}
}

func envListCmd() *cli.Command {
	args := workflowArgs
	args.Validator = cli.ValidateFunc(func(args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("expects at most 1 arg, received %v", len(args))
		}
		return nil
	})

	cmd := &cli.Command{
		Use:     "list [<path>]",
		Aliases: []string{"ls"},
		Short:   "list environments relative to current dir or <path>",
		Args:    args,
	}

	useJSON := cmd.Flags().Bool("json", false, "json output")
	getLabelSelector := labelSelectorFlag(cmd.Flags())

	useNames := cmd.Flags().Bool("names", false, "plain names output")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		var path string
		var err error
		if len(args) == 1 {
			path = args[0]
		} else {
			path, err = os.Getwd()
			if err != nil {
				return nil
			}
		}

		envs, err := tanka.ListEnvs(path, tanka.ListOpts{Selector: getLabelSelector()})
		if err != nil {
			return err
		}

		if *useJSON {
			j, err := json.Marshal(envs)
			if err != nil {
				return fmt.Errorf("Formatting as json: %s", err)
			}
			fmt.Println(string(j))
			return nil
		} else if *useNames {
			for _, e := range envs {
				fmt.Println(e.Metadata.Name)
			}
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		f := "%s\t%s\t%s\t\n"
		fmt.Fprintf(w, f, "NAME", "NAMESPACE", "SERVER")
		for _, e := range envs {
			fmt.Fprintf(w, f, e.Metadata.Name, e.Spec.Namespace, e.Spec.APIServer)
		}
		w.Flush()

		return nil
	}
	return cmd
}
