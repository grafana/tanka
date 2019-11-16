package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/fatih/structs"
	"github.com/spf13/cobra"

	"github.com/grafana/tanka/pkg/tanka"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <path>",
		Short: "display an overview",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{
			"args": "baseDir",
		},
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		status, err := tanka.Status(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Context:", status.Client.Context.Get("name"))
		fmt.Println("Cluster:", status.Client.Context.Get("context").MustObjxMap().Get("cluster"))
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
	}
	return cmd
}
