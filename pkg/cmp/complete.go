package cmp

import (
	"github.com/posener/complete"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CompletionHandlers is a set of predictors
type CompletionHandlers map[string]complete.Predictor

// Add adds a predictor to the list of handlers
func (h CompletionHandlers) Add(k string, p complete.Predictor) {
	h[k] = p
}

// Get returns a predictor from the list of handlers
func (h CompletionHandlers) Get(k string) complete.Predictor {
	return h[k]
}

// GetOrNone returns a predictor in case it exists, otherwise complete.PredictNothing
func (h CompletionHandlers) GetOrNone(k string) complete.Predictor {
	if p, ok := h[k]; ok {
		return p
	}
	return complete.PredictNothing
}

// Has returns whether a predictor exists
func (h CompletionHandlers) Has(k string) bool {
	_, ok := h[k]
	return ok
}

// Handlers are global Handlers to be used in annotations
var Handlers = CompletionHandlers{}

func init() {
	Handlers["dirs"] = complete.PredictDirs("*")
}

// Create parses a *cobra.Command into a complete.Command
func Create(root *cobra.Command) complete.Command {
	rootCmp := complete.Command{}

	// Flags
	rootCmp.Flags = make(complete.Flags)
	addFlags := func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		handler := Handlers.GetOrNone(root.Annotations["flags/"+flag.Name])

		if len(flag.Shorthand) > 0 {
			rootCmp.Flags["-"+flag.Shorthand] = handler
		}

		rootCmp.Flags["--"+flag.Name] = handler
	}
	root.LocalFlags().VisitAll(addFlags)
	root.InheritedFlags().VisitAll(addFlags)

	// Subcommands
	if root.HasAvailableSubCommands() {
		rootCmp.Sub = make(complete.Commands)
		for _, c := range root.Commands() {
			if !c.Hidden {
				rootCmp.Sub[c.Name()] = Create(c)
			}
		}
	}

	// Positional Arguments
	rootCmp.Args = Handlers.GetOrNone(root.Annotations["args"])
	return rootCmp
}
