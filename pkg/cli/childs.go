package cli

import "fmt"

// AddChildren adds the supplied commands as subcommands.
// This command is set as the parent of the new children.
func (c *Command) AddCommand(childs ...*Command) {
	if c.children == nil {
		c.children = make(map[string]*Command)
	}

	for _, child := range childs {
		child.parentPtr = c
		c.children[child.Name()] = child
	}
}

// findTarget finds the specified (sub)command based on the args
func findTarget(c *Command, args []string) (*Command, []string, error) {
	argsWOflags := stripFlags(args, c)
	if len(argsWOflags) == 0 {
		return c, args, nil
	}
	nextSubCmd := argsWOflags[0]

	cmd, ok := c.child(nextSubCmd)
	switch {
	case !ok:
		return nil, nil, c.help(fmt.Errorf("unknown subcommand `%s`", nextSubCmd))
	case cmd != nil:
		return findTarget(cmd, argsMinusFirstX(args, nextSubCmd))
	}
	return c, args, nil
}

func (c *Command) child(name string) (*Command, bool) {
	child, ok := c.children[name]
	return child, ok
}

// argsMinusFirstX removes only the first x from args.  Otherwise, commands that look like
// openshift admin policy add-role-to-user admin my-user, lose the admin argument (arg[4]).
func argsMinusFirstX(args []string, x string) []string {
	for i, y := range args {
		if x == y {
			ret := []string{}
			ret = append(ret, args[:i]...)
			ret = append(ret, args[i+1:]...)
			return ret
		}
	}
	return args
}
