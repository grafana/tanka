package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

// initCmd creates a new application
func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create the directory structure",
		Args:  cobra.NoArgs,
	}

	force := cmd.Flags().BoolP("force", "f", false, "ignore the working directory not being empty")
	installK8sLibFlag := cmd.Flags().Bool("k8s", true, "set to false to skip installation of k.libsonnet")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		failed := false

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

		if *installK8sLibFlag {
			if err := installK8sLib(); err != nil {
				// This is not fatal, as most of Tanka will work anyways
				log.Println("Installing k.libsonnet:", err)
				failed = true
			}
		}

		fmt.Println("Directory structure set up! Remember to configure the API endpoint:\n`tk env set environments/default --server=127.0.0.1:6443`")
		if failed {
			log.Println("Errors occured while initializing the project. Check the above logs for details.")
		}
	}
	return cmd
}

func installK8sLib() error {
	if _, err := exec.LookPath("jb"); err != nil {
		return errors.New("jsonnet-bundler not found in $PATH. Follow https://tanka.dev/install#jsonnet-bundler for installation instructions.")
	}

	// TODO: use the jb packages for this once refactored there
	const klibsonnetJsonnetfile = `{
"dependencies": [
    {
      "source": {
        "git": {
          "remote": "https://github.com/grafana/jsonnet-libs",
          "subdir": "ksonnet-util"
        }
      },
      "version": "master"
    },
    {
      "source": {
        "git": {
          "remote": "https://github.com/ksonnet/ksonnet-lib",
          "subdir": "ksonnet.beta.4"
        }
      },
      "version": "master"
    }
  ]
}
`

	if err := writeNewFile("lib/k.libsonnet", `import "ksonnet.beta.4/k.libsonnet"`); err != nil {
		return err
	}

	if err := ioutil.WriteFile("jsonnetfile.json", []byte(klibsonnetJsonnetfile), 0644); err != nil {
		return err
	}

	cmd := exec.Command("jb", "install")
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
