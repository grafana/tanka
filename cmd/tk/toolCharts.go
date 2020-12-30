package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-clix/cli"
	"github.com/grafana/tanka/pkg/helm"
	"gopkg.in/yaml.v2"
)

func chartsCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "charts",
		Short: "Declarative vendoring of Helm Charts",
	}

	cmd.AddCommand(
		chartsInitCmd(),
		chartsAddCmd(),
		chartsAddRepoCmd(),
		chartsVendorCmd(),
		chartsConfigCmd(),
	)

	return cmd
}

func chartsVendorCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "vendor",
		Short: "Download Charts to a local folder",
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		c, err := loadChartfile()
		if err != nil {
			return err
		}

		return c.Vendor()
	}

	return cmd
}

func chartsAddCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "add [chart@version] [...]",
		Short: "Adds Charts to the chartfile",
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		c, err := loadChartfile()
		if err != nil {
			return err
		}

		return c.Add(args)
	}

	return cmd
}

func chartsAddRepoCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "add-repo [NAME] [URL]",
		Short: "Adds a repository to the chartfile",
		Args:  cli.ArgsExact(2),
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		c, err := loadChartfile()
		if err != nil {
			return err
		}

		return c.AddRepos(helm.Repo{
			Name: args[0],
			URL:  args[1],
		})
	}

	return cmd
}

func chartsConfigCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "config",
		Short: "Displays the current manifest",
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		c, err := loadChartfile()
		if err != nil {
			return err
		}

		data, err := yaml.Marshal(c.Manifest)
		if err != nil {
			return err
		}

		fmt.Print(string(data))

		return nil
	}

	return cmd
}

func chartsInitCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "init",
		Short: "Create a new Chartfile",
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		path := filepath.Join(wd, helm.Filename)
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("Chartfile at '%s' already exists. Aborting", path)
		}

		if _, err := helm.InitChartfile(path); err != nil {
			return err
		}

		log.Printf("Success! New Chartfile created at '%s'", path)
		return nil
	}

	return cmd
}

func loadChartfile() (*helm.Charts, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return helm.LoadChartfile(wd)
}
