package main

import (
	"github.com/grafana/tanka/pkg/cli"
	"github.com/posener/complete"
)

var workflowArgs = cli.Args{
	Validator: cli.ValidateExact(1),
	Predictor: cli.PredictFunc(func(complete.Args) []string {
		if dirs := findBaseDirs(); len(dirs) != 0 {
			return dirs
		}
		return []string{""}
	}),
}
