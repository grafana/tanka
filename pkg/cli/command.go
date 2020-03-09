package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/posener/complete"
	"github.com/spf13/pflag"
)

// Command represents a (sub)command of the application. Either `Run()` must be
// defined, or subcommands added using `AddCommand()`. These are also mutually
// exclusive.
type Command struct {
	// Usage line. First word must be the command name, everything else is
	// displayed as-is.
	Use string
	// Aliases define alternative names for a command
	Aliases []string

	// Short help text, used for overviews
	Short string
	// Long help text, used for full help pages. `Short` is used as a fallback
	// if unset.
	Long string

	// Version of the application. Only used on the root command
	Version string

	// Run is the action that is run when this command is invoked.
	// The error is returned as-is from `Execute()`.
	Run func(cmd *Command, args []string) error

	// Validation + Completion
	//
	// Predict contains Predictors for flags. Defaults to
	// `complete.PredictSomething` if unset.
	// Use the flags name (not shorthand) as the key.
	Predictors map[string]complete.Predictor
	// Args is used to validate and complete positional arguments
	Args Arguments

	// internal fields
	children  []*Command
	flags     *pflag.FlagSet
	parentPtr *Command
}

// Execute runs the application. It should be run on the most outer level
// command.
// The error return value is used for both, application errors but also help texts.
func (c *Command) Execute() error {
	// Execute must be called on the top level command
	if c.parentPtr != nil {
		return c.parentPtr.Execute()
	}

	// add subcommand for install CLI completions
	c.AddCommand(completionCmd(c.Use))

	// add version flag to root command. If `findTarget` switches to a
	// subcommand afterwards, showVersion will stay nil, causing the version not
	// to be shown.
	var showVersion *bool
	if c.Version != "" {
		showVersion = c.Flags().Bool("version", false, fmt.Sprintf("version for %s", c.Use))
	}

	// exit if in bash completion context
	if predict(c) {
		return nil
	}

	// find the correct (sub)command
	c, args, err := findTarget(c, os.Args[1:])
	if err != nil {
		return err
	}

	// add help flag
	var showHelp *bool
	if c.Flags().Lookup("help") == nil {
		showHelp = initHelpFlag(c)
	}

	// parse flags
	if err := c.Flags().Parse(args); err != nil {
		return c.help(err)
	}

	// show version if requested. This only happens if findTarget returned the
	// root command.
	if showVersion != nil && *showVersion {
		log.Printf("%s version %s", c.Use, c.Version)
		return nil
	}

	// show help if requested or missing `Run()`
	switch {
	case showHelp != nil && *showHelp:
		fallthrough
	case c.Run == nil:
		return helpErr(c)
	}

	// validate args
	if err := c.Args.Validate(c.Flags().Args()); err != nil {
		return c.help(err)
	}

	// run!
	return c.Run(c, c.Flags().Args())
}

func initHelpFlag(c *Command) *bool {
	return c.Flags().BoolP("help", "h", false, "help for "+c.Name())
}

func helpErr(c *Command) error {
	help := c.Short
	if c.Long != "" {
		help = c.Long
	}

	return fmt.Errorf("%s\n\n%s", help, c.Usage())
}

// Name of this command. The first segment of the `Use` field.
func (c *Command) Name() string {
	return strings.Split(c.Use, " ")[0]
}

// Usage string
func (c *Command) Usage() string {
	return c.helpable().Generate()
}

func (c *Command) helpable() *helpable {
	return &helpable{*c}
}
