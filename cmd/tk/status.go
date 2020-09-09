package main

import (
	"fmt"
	"os"
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
		fmt.Println("Environment:")
		for k, v := range structs.Map(status.Env.Spec) {
			fmt.Printf("  %s: %s\n", k, v)
		}

		fmt.Println("Resources:")
		f := "  %s\t%s/%s\n"
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "  NAMESPACE\tOBJECTSPEC")
		for _, r := range status.Resources {
			fmt.Fprintf(w, f, r.Metadata().Namespace(), r.Kind(), r.Metadata().Name())
		}
		w.Flush()

		return nil
	}
	return cmd
}
