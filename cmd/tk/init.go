package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sh0rez/tanka/pkg/config/v1alpha1"
	"github.com/spf13/cobra"
)

// initCmd creates a new application
func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create the directory structure",
	}
	force := cmd.Flags().BoolP("force", "f", false, "ignore the working directory not being empty")
	provider := cmd.Flags().StringP("provider", "p", "", fmt.Sprintf("one of %s", listProviders()))
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

		if err := os.MkdirAll("environments/default", os.ModePerm); err != nil {
			log.Fatalln("Error creating environments folder")
		}

		if err := writeNewFile("environments/default/main.jsonnet", "{}"); err != nil {
			log.Fatalln("Error creating `main.jsonnet`:", err)
		}

		cfg := v1alpha1.Config{
			APIVersion: "tanka.dev/v1alpha1",
			Kind:       "Environment",
		}
		if *provider != "" {
			cfg.Spec = map[string]interface{}{}
			cfg.Spec["provider"] = map[string]interface{}{
				*provider: map[string]interface{}{},
			}
		}

		spec, err := json.Marshal(&cfg)
		if err != nil {
			log.Fatalln("Error creating spec.json:", err)
		}

		if err := writeNewFile("environments/default/spec.json", string(spec)); err != nil {
			log.Fatalln("Error creating `environments/default/spec.json`:", err)
		}

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
