package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/grafana/tanka/pkg/config/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func envCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env [action]",
		Short: "manipulate environments",
	}
	cmd.PersistentFlags().Bool("json", false, "output in json format")
	cmd.AddCommand(envListCmd())
	return cmd
}

func envListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list environments",
		Args:  cobra.NoArgs,
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		envs := []v1alpha1.Config{}
		dirs := findBaseDirs()
		useJson, err := cmd.Flags().GetBool("json")
		if err != nil {
			// this err should never occur. Panic in case
			panic(err)
		}
		for _, dir := range dirs {
			viper.Reset()
			envs = append(envs, *setupConfiguration(dir))
		}

		if useJson {
			j, err := json.Marshal(envs)
			if err != nil {
				log.Fatalln("Formatting as json:", j)
			}
			fmt.Println(string(j))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		f := "%s\t%s\t%s\t\n"
		fmt.Fprintf(w, f, "NAME", "NAMESPACE", "SERVER")
		for _, e := range envs {
			fmt.Fprintf(w, f, e.Metadata.Name, e.Spec.Namespace, e.Spec.APIServer)
		}
		w.Flush()
	}
	return cmd
}
