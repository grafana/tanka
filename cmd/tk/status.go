package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/structs"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/tanka"
)

func statusCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "status <path>",
		Short: "display an overview of the environment, including contents and metadata.",
		Args:  workflowArgs,
	}

	getJsonnetOpts := jsonnetFlags(cmd.Flags())

	cmd.Run = func(cmd *cli.Command, args []string) error {
		status, err := tanka.Status(args[0], tanka.Opts{
			JsonnetOpts: getJsonnetOpts(),
		})
		if err != nil {
			return err
		}

		context := status.Client.Kubeconfig.Context
		fmt.Println("Context:", context.Name)
		fmt.Println("Cluster:", context.Context.Cluster)

		fmt.Println()
		fmt.Println("Directory Layout:")
		fmt.Println("  Project root:", status.DirLayout.Root)
		fmt.Println("  Environment:", status.DirLayout.Base)
		fmt.Println("  Entrypoint:", status.DirLayout.Entrypoint)

		fmt.Println()
		fmt.Println("Environment:")
		for _, f := range structs.Fields(status.Env.Spec) {
			fmt.Printf("  %s: %v\n", f.Name(), mustJsonFmt(f.Value()))
		}

		fmt.Println()
		fmt.Println("Resources:")
		f := "  %s\t%s/%s\n"
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "  NAMESPACE\t<KIND>/<NAME>")
		for _, r := range status.Resources {
			fmt.Fprintf(w, f, r.Metadata().Namespace(), r.Kind(), r.Metadata().Name())
		}
		w.Flush()

		return nil
	}
	return cmd
}

func mustJsonFmt(i interface{}) string {
	data, err := json.Marshal(i)
	if err != nil {
		return err.Error()
	}

	out, err := tanka.Format("", string(data))
	if err != nil {
		return err.Error()
	}

	return strings.TrimSuffix(out, "\n")
}
