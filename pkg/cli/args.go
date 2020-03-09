package cli

import (
	"fmt"

	"github.com/posener/complete"
)

// Arguments is used to validate and complete positional arguments.
// Use `Args()` to create an instance from functions.
type Arguments interface {
	Validator
	complete.Predictor
}

type Validator interface {
	// Validate receives the arguments of the command (without flags) and shall
	// return an error if they are unexpected.
	Validate(args []string) error
}

type Args struct {
	Validator
	complete.Predictor
}

type ValidateFunc func(args []string) error

func (v ValidateFunc) Validate(args []string) error {
	return v(args)
}

type PredictFunc = complete.PredictFunc

// No Arguments
func ValidateNone() ValidateFunc {
	return ValidateExact(0)
}

func PredictNone() complete.Predictor {
	return complete.PredictNothing
}

func ArgsNone() Arguments {
	return Args{
		Validator: ValidateNone(),
		Predictor: PredictNone(),
	}
}

// Exact arguments
func ValidateExact(n int) ValidateFunc {
	return func(args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %v arg, received %v", n, len(args))
		}
		return nil
	}
}

func ArgsExact(n int) Arguments {
	return Args{
		Validator: ValidateExact(n),
		Predictor: PredictAny(),
	}
}

// Any arguments
func PredictAny() complete.Predictor {
	return complete.PredictAnything
}

func ValidateAny() ValidateFunc {
	return func(args []string) error {
		return nil
	}
}

func ArgsAny() Arguments {
	return Args{
		Validator: ValidateAny(),
		Predictor: PredictAny(),
	}
}

// Predefined arguments
func ValidateSet(set ...string) ValidateFunc {
	return func(args []string) error {
		if err := ValidateExact(1)(args); err != nil {
			return err
		}

		for _, s := range set {
			if args[0] == s {
				return nil
			}
		}

		return fmt.Errorf("only accepts %v", set)
	}
}

func PredictSet(set ...string) complete.Predictor {
	return complete.PredictSet(set...)
}

func ArgsSet(set ...string) Arguments {
	return Args{
		Validator: ValidateSet(set...),
		Predictor: PredictSet(set...),
	}
}
