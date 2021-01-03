package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-clix/cli"
	"github.com/posener/complete"

	"github.com/grafana/tanka/pkg/tanka"
)

// ValidateMin checks that at least n arguments were given
func ValidateMin(n int) cli.ValidateFunc {
	return func(args []string) error {
		if len(args) < n {
			return fmt.Errorf("expects at least %v arg, received %v", n, len(args))
		}
		return nil
	}
}

var workflowArgs = cli.Args{
	Validator: ValidateMin(1),
	Predictor: cli.PredictFunc(func(args complete.Args) []string {
		pwd, err := os.Getwd()
		if err != nil {
			return nil
		}

		dirs, err := tanka.FindBaseDirs(pwd)
		if err != nil {
			return nil
		}

		var reldirs []string
		for _, dir := range dirs {
			reldir, err := filepath.Rel(pwd, dir)
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
