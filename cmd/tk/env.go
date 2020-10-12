package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/go-clix/cli"
	"github.com/pkg/errors"
	"github.com/posener/complete"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/grafana/tanka/pkg/kubernetes/client"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
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

func envSettingsFlags(env *v1alpha1.Config, fs *pflag.FlagSet) {
	fs.StringVar(&env.Spec.APIServer, "server", env.Spec.APIServer, "endpoint of the Kubernetes API")
	fs.StringVar(&env.Spec.APIServer, "server-from-context", env.Spec.APIServer, "set the server to a known one from $KUBECONFIG")
	fs.StringVar(&env.Spec.Namespace, "namespace", env.Spec.Namespace, "namespace to create objects in")
	fs.StringVar(&env.Spec.DiffStrategy, "diff-strategy", env.Spec.DiffStrategy, "specify diff-strategy. Automatically detected otherwise.")
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
	tmp := v1alpha1.Config{}
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

		cfg := setupConfiguration(path)
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
	cmd.Run = func(cmd *cli.Command, args []string) error {
		if cmd.Flags().Changed("server-from-context") {
			server, err := client.IPFromContext(cfg.Spec.APIServer)
			if err != nil {
				return fmt.Errorf("Resolving IP from context: %s", err)
			}
			cfg.Spec.APIServer = server
		}

		if err := addEnv(args[0], cfg); err != nil {
			return err
		}

		return nil
	}
	return cmd
}

// used by initCmd() as well
func addEnv(dir string, cfg *v1alpha1.Config) error {
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

	// write spec.json
	if err := writeJSON(cfg, filepath.Join(path, "spec.json")); err != nil {
		return err
	}

	// write main.jsonnet
	if err := writeJSON(struct{}{}, filepath.Join(path, "main.jsonnet")); err != nil {
		return err
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
	cmd := &cli.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list environments",
		Args:    cli.ArgsNone(),
	}

	useJSON := cmd.Flags().Bool("json", false, "json output")
	labelSelector := cmd.Flags().StringP("selector", "l", "", "Label selector. Uses the same syntax as kubectl does")

	useNames := cmd.Flags().Bool("names", false, "plain names output")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		envs := []v1alpha1.Config{}
		dirs := findBaseDirs()
		var selector labels.Selector
		var err error

		if *labelSelector != "" {
			selector, err = labels.Parse(*labelSelector)
			if err != nil {
				return err
			}
		}

		for _, dir := range dirs {
			env := setupConfiguration(dir)
			if env == nil {
				log.Printf("Could not setup configuration from %q", dir)
				continue
			}
			if selector == nil || selector.Empty() || selector.Matches(env.Metadata) {
				envs = append(envs, *env)
			}
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
