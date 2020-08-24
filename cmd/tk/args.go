package main

import (
	"github.com/go-clix/cli"
	"github.com/posener/complete"
)

var workflowArgs = cli.Args{
	Validator: cli.ValidateExact(1),
	Predictor: cli.PredictFunc(func(args complete.Args) []string {
		if dirs := findBaseDirs(); len(dirs) != 0 {
			return dirs
		}

		return complete.PredictDirs("*").Predict(args)
	}),
}
