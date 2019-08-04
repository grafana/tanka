package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh0rez/tanka/pkg/jpath"
	"github.com/sh0rez/tanka/pkg/jsonnet"
	"github.com/spf13/cobra"
)

func fmtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "format .jsonnet and .libsonnet files",
		Use:   "fmt",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {}
	return cmd
}

func evalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "evaluate the jsonnet to json",
		Use:   "eval",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		json, err := eval(cmd.Flag("file").Value.String())
		if err != nil {
			return err
		}
		fmt.Print(json)
		return nil
	}

	return cmd
}

func eval(filename string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	_, baseDir, _ := jpath.Resolve(pwd, filename)
	json, err := jsonnet.EvaluateFile(filepath.Join(baseDir, filename))
	if err != nil {
		return "", err
	}
	return json, nil
}

func evalDict(filename string) (map[string]interface{}, error) {
	var rawDict map[string]interface{}

	raw, err := eval(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(raw), &rawDict); err != nil {
		return nil, err
	}
	return rawDict, nil
}
