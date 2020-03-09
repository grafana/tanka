package cli

import (
	"github.com/posener/complete"
	"github.com/spf13/pflag"
)

func predict(c *Command) bool {
	if c.Args == nil {
		c.Args = ArgsAny()
	}

	cmp := complete.New(c.Name(), createCmp(c))
	return cmp.Complete()
}

func createCmp(c *Command) complete.Command {
	rootCmp := complete.Command{}

	rootCmp.Flags = complete.Flags{
		"-h":     complete.PredictNothing,
		"--help": complete.PredictNothing,
	}

	c.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		p, ok := c.Predictors[flag.Name]
		if !ok {
			p = complete.PredictNothing
		}

		if len(flag.Shorthand) > 0 {
			rootCmp.Flags["-"+flag.Shorthand] = p
		}
		rootCmp.Flags["--"+flag.Name] = p
	})

	if c.children != nil {
		rootCmp.Sub = make(complete.Commands)
		for _, c := range c.children {
			rootCmp.Sub[c.Name()] = createCmp(c)
		}
	}

	// Positional Arguments
	rootCmp.Args = c.Args
	if rootCmp.Args == nil {
		rootCmp.Args = PredictAny()
	}

	return rootCmp
}
