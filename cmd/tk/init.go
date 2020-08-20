package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// initCmd creates a new application
func initCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "init",
		Short: "Create the directory structure",
		Args:  cli.ArgsNone(),
	}

	force := cmd.Flags().BoolP("force", "f", false, "ignore the working directory not being empty")
	installK8sLibFlag := cmd.Flags().Bool("k8s", true, "set to false to skip installation of k.libsonnet")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		failed := false

		files, err := ioutil.ReadDir(".")
		if err != nil {
			return fmt.Errorf("Error listing files: %s", err)
		}
		if len(files) > 0 && !*force {
			return fmt.Errorf("Error: directory not empty. Use `-f` to force")
		}

		if err := writeNewFile("jsonnetfile.json", "{}"); err != nil {
			return fmt.Errorf("Error creating `jsonnetfile.json`: %s", err)
		}

		if err := os.Mkdir("vendor", os.ModePerm); err != nil {
			return fmt.Errorf("Error creating `vendor/` folder: %s", err)
		}

		if err := os.Mkdir("lib", os.ModePerm); err != nil {
			return fmt.Errorf("Error creating `lib/` folder: %s", err)
		}

		cfg := v1alpha1.New()
		if err := addEnv("environments/default", cfg); err != nil {
			return err
		}

		if *installK8sLibFlag {
			if err := installK8sLib(); err != nil {
				// This is not fatal, as most of Tanka will work anyways
				log.Println("Installing k.libsonnet:", err)
				failed = true
			}
		}

		fmt.Println("Directory structure set up! Remember to configure the API endpoint:\n`tk env set environments/default --server=https://127.0.0.1:6443`")
		if failed {
			log.Println("Errors occured while initializing the project. Check the above logs for details.")
		}

		return nil
	}
	return cmd
}

func installK8sLib() error {
	jbBinary := "jb"
	if env := os.Getenv("TANKA_JB_PATH"); env != "" {
		jbBinary = env
	}

	if _, err := exec.LookPath(jbBinary); err != nil {
		return errors.New("jsonnet-bundler not found in $PATH. Follow https://tanka.dev/install#jsonnet-bundler for installation instructions")
	}

	var initialPackages = []string{
		"github.com/ksonnet/ksonnet-lib/ksonnet.beta.4",
		"github.com/grafana/jsonnet-libs/ksonnet-util",
	}

	if err := writeNewFile("lib/k.libsonnet", "import 'ksonnet.beta.4/k.libsonnet'\n"); err != nil {
		return err
	}

	cmd := exec.Command(jbBinary, append([]string{"install"}, initialPackages...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// writeNewFile writes the content to a file if it does not exist
func writeNewFile(name, content string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return ioutil.WriteFile(name, []byte(content), 0644)
	}
	return nil
}
