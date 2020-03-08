package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// Command represents a (sub)command of the application. Either `Run()` must be
// defined, or subcommands added using `AddCommand()`. These are also mutually
// exclusive.
type Command struct {
	// Usage line. First word must be the command name, everything else is
	// displayed as-is.
	Use string

	// Short help text, used for overviews
	Short string
	// Long help text, used for full help pages. `Short` is used as a fallback
	// if unset.
	Long string

	// Run is the action that is run when this command is invoked.
	// The error is returned as-is from `Execute()`.
	Run func(cmd *Command, args []string) error

	// internal fields
	children  map[string]*Command
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

	c, args, err := c.find(os.Args[1:])
	if err != nil {
		return err
	}

	// add help flag
	var showHelp *bool
	if c.Flags().Lookup("help") == nil {
		showHelp = c.Flags().BoolP("help", "h", false, "help for "+c.Name())
	}

	if err := c.Flags().Parse(args); err != nil {
		return c.help(err)
	}

	switch {
	case showHelp != nil && *showHelp:
		fallthrough
	case c.Run == nil:
		return helpErr(c)
	}

	return c.Run(c, c.Flags().Args())
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
