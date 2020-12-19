package main

import (
	"os"

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

		dirs, err := tanka.FindBaseDirs(pwd)
		if err != nil {
			return nil
		}

		if len(dirs) != 0 {
			return dirs
		}

		return complete.PredictDirs("*").Predict(args)
	}),
}
