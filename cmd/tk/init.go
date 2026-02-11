package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-clix/cli"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

const defaultK8sVersion = "1.32"

// Pinned commit of k8s-libsonnet that contains defaultK8sVersion.
// This avoids failures when older versions are removed from main branch.
// See https://github.com/grafana/tanka/issues/1863
const k8sLibsonnetCommitRef = "55380470fb7979e6ce0c4316cb9c27a266caf298"

// initCmd creates a new application
func initCmd(ctx context.Context) *cli.Command {
	cmd := &cli.Command{
		Use:   "init",
		Short: "Create the directory structure",
		Args:  cli.ArgsNone(),
	}

	force := cmd.Flags().BoolP("force", "f", false, "ignore the working directory not being empty")
	installK8s := cmd.Flags().String("k8s", defaultK8sVersion, "choose the version of k8s-libsonnet, full package URI, or false to skip (e.g. \"1.32\", \"github.com/jsonnet-libs/k8s-libsonnet/1.32@main\")")
	inline := cmd.Flags().BoolP("inline", "i", false, "create an inline environment")

	cmd.Run = func(_ *cli.Command, _ []string) error {
		_, span := tracer.Start(ctx, "initCmd")
		defer span.End()
		failed := false

		files, err := os.ReadDir(".")
		if err != nil {
			return fmt.Errorf("error listing files: %s", err)
		}
		if len(files) > 0 && !*force {
			return fmt.Errorf("error: directory not empty. Use `-f` to force")
		}

		if err := writeNewFile("jsonnetfile.json", "{}"); err != nil {
			return fmt.Errorf("error creating `jsonnetfile.json`: %s", err)
		}

		if err := os.Mkdir("vendor", os.ModePerm); err != nil {
			return fmt.Errorf("error creating `vendor/` folder: %s", err)
		}

		if err := os.Mkdir("lib", os.ModePerm); err != nil {
			return fmt.Errorf("error creating `lib/` folder: %s", err)
		}

		cfg := v1alpha1.New()
		if err := addEnv("environments/default", cfg, *inline); err != nil {
			return err
		}

		version := *installK8s
		doInstall, err := strconv.ParseBool(*installK8s)
		if err != nil {
			// --k8s=<non-boolean>
			doInstall = true
		} else {
			// --k8s=<boolean>, fallback to default version
			version = defaultK8sVersion
		}

		if doInstall {
			if err := installK8sLib(version); err != nil {
				// This is not fatal, as most of Tanka will work anyways
				fmt.Println("Installing k.libsonnet:", err)
				failed = true
			}
		}

		if *inline {
			fmt.Println("Directory structure set up! Remember to configure the API endpoint in environments/default/main.jsonnet")
		} else {
			fmt.Println("Directory structure set up! Remember to configure the API endpoint:\n`tk env set environments/default --server=https://127.0.0.1:6443`")
		}
		if failed {
			fmt.Println("Errors occurred while initializing the project. Check the above logs for details.")
		}

		return nil
	}
	return cmd
}

// The version can be:
// - a full package URI (e.g. "github.com/jsonnet-libs/k8s-libsonnet/1.32@main")
// - a version number (e.g. "1.32"): use default package with pinned commit
func installK8sLib(version string) error {
	jbBinary := "jb"
	if env := os.Getenv("TANKA_JB_PATH"); env != "" {
		jbBinary = env
	}

	if _, err := exec.LookPath(jbBinary); err != nil {
		return errors.New("jsonnet-bundler not found in $PATH. Follow https://tanka.dev/install#jsonnet-bundler for installation instructions")
	}

	k8sLibsonnetURI := version
	if !strings.Contains(k8sLibsonnetURI, "/") {
		// If it doesn't look like a full package URI, it's a version number.
		k8sLibsonnetURI = "github.com/jsonnet-libs/k8s-libsonnet/" + version
	}
	importPathPrefix, _, ok := strings.Cut(k8sLibsonnetURI, "@")
	if !ok {
		k8sLibsonnetURI += "@" + k8sLibsonnetCommitRef
	}

	var initialPackages = []string{
		k8sLibsonnetURI,
		"github.com/grafana/jsonnet-libs/ksonnet-util",
		"github.com/jsonnet-libs/docsonnet/doc-util", // install docsonnet to make `tk lint` work
	}

	if err := writeNewFile("lib/k.libsonnet", "import '"+importPathPrefix+"/main.libsonnet'\n"); err != nil {
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
		return os.WriteFile(name, []byte(content), 0644)
	}
	return nil
}
