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

// Validator checks that arguments have the expected form
type Validator interface {
	// Validate receives the arguments of the command (without flags) and shall
	// return an error if they are unexpected.
	Validate(args []string) error
}

// Args bundles user-supplied implementations of the respective interfaces into
// an Arguments implementation.
type Args struct {
	Validator
	complete.Predictor
}

// ValidateFunc allows to use an ordinary func as an Validator
type ValidateFunc func(args []string) error

// Validate wrap the underlying function
func (v ValidateFunc) Validate(args []string) error {
	return v(args)
}

// PredictFunc allows to use an ordinary func as an Predictor
type PredictFunc = complete.PredictFunc

// ---
// Common Argument implementations
// ---

// No Arguments

// ArgsNone checks that no arguments were given, and disables predictions.
func ArgsNone() Arguments {
	return Args{
		Validator: ValidateNone(),
		Predictor: PredictNone(),
	}
}

// ValidateNone checks for no arguments at all
func ValidateNone() ValidateFunc {
	return ValidateExact(0)
}

// PredictNone predicts exactly nothing
func PredictNone() complete.Predictor {
	return PredictFunc(func(args complete.Args) []string {
		return nil
	})
}

// Exact arguments

// ArgsExact checks for exactly n arguments, predicting anything
func ArgsExact(n int) Arguments {
	return Args{
		Validator: ValidateExact(n),
		Predictor: PredictAny(),
	}
}

// ValidateExact checks that exactly n arguments were given
func ValidateExact(n int) ValidateFunc {
	return func(args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %v arg, received %v", n, len(args))
		}
		return nil
	}
}

// Any arguments

// ArgsAny allows any number of arguments with any value
func ArgsAny() Arguments {
	return Args{
		Validator: ValidateAny(),
		Predictor: PredictAny(),
	}
}

// PredictAny predicts any files/directories
func PredictAny() complete.Predictor {
	return complete.PredictFiles("*")
}

// ValidateAny always approves
func ValidateAny() ValidateFunc {
	return func(args []string) error {
		return nil
	}
}

// Predefined arguments

// ArgsSet check the given argument is in the predefined set of options. Only
// these options are predicted. Only a single argument is assumed.
func ArgsSet(set ...string) Arguments {
	return Args{
		Validator: ValidateSet(set...),
		Predictor: PredictSet(set...),
	}
}

// PredictSet predicts the values from the given set
func PredictSet(set ...string) complete.Predictor {
	return complete.PredictSet(set...)
}

// ValidateSet checks that the given single argument is part of the set.
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
