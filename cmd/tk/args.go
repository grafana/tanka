package main

import (
	"os"
	"path/filepath"

	"github.com/go-clix/cli"
	"github.com/posener/complete"

	"github.com/grafana/tanka/pkg/tanka"
)

var workflowArgs = cli.Args{
	Validator: cli.ValidateExact(1),
	Predictor: cli.PredictFunc(func(args complete.Args) []string {
		pwd, err := os.Getwd()
		if err != nil {
			return nil
		}

		dirs, err := tanka.ListEnvs(pwd, tanka.ListOpts{})
		if err != nil {
			return nil
		}

		var reldirs []string
		for _, env := range dirs {
			reldir, err := filepath.Rel(pwd, env.Metadata.Namespace)
			if err == nil {
				reldirs = append(reldirs, reldir)
			}
		}

		if len(reldirs) != 0 {
			return reldirs
		}

		return complete.PredictDirs("*").Predict(args)
	}),
}
