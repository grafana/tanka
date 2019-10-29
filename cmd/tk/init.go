package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/spf13/cobra"
)

// initCmd creates a new application
func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create the directory structure",
		Args:  cobra.NoArgs,
	}
	force := cmd.Flags().BoolP("force", "f", false, "ignore the working directory not being empty")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatalln("Error listing files:", err)
		}
		if len(files) > 0 && !*force {
			log.Fatalln("Error: directory not empty. Use `-f` to force")
		}

		if err := writeNewFile("jsonnetfile.json", "{}"); err != nil {
			log.Fatalln("Error creating `jsonnetfile.json`:", err)
		}

		if err := os.Mkdir("vendor", os.ModePerm); err != nil {
			log.Fatalln("Error creating `vendor/` folder:", err)
		}

		if err := os.Mkdir("lib", os.ModePerm); err != nil {
			log.Fatalln("Error creating `vendor/` folder:", err)
		}

		cfg := v1alpha1.New()
		if err := addEnv("environments/default", cfg); err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Directory structure set up! Remember to configure the API endpoint:\n`tk env set environments/default --server=127.0.0.1:6443`")
	}
	return cmd
}

// writeNewFile writes the content to a file if it does not exist
func writeNewFile(name, content string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return ioutil.WriteFile(name, []byte(content), 0644)
	}
	return nil
}
