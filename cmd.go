// Package cmd implements a command-line parser.
package cmd

import (
	"errors"
	"fmt"
	"os"
)

// A Cmd represents a command with command-line flags and positional arguments.
type Cmd struct {
	Flags
	Summary, Details      string
	name                  string
	args                  []arg
	hasOptional, hasMulti bool
}

type arg struct {
	name     string
	optional bool
	single   *string
	multi    *[]string
}

// New returns a new command with the specified name.
func New(name string) *Cmd {
	return nil
}

// Arg defines a positional argument.
func (c *Cmd) Arg(name string, p *string) {
	if c.hasOptional {
		panic("Cmd: cannot add non-optional after optional argument")
	}
	if c.hasMulti {
		panic("Cmd: cannot add single argument after repeated argument")
	}
	c.args = append(c.args, arg{
		name:   name,
		single: p,
	})
}

// OptionalArg defines an optional positional argument.
func (c *Cmd) OptionalArg(name string, p *string) {
	if c.hasMulti {
		panic("Cmd: cannot add single argument after repeated argument")
	}
	c.args = append(c.args, arg{
		name:     name,
		optional: true,
		single:   p,
	})
	c.hasOptional = true
}

// Args defines an argument that can be present one or more times.
func (c *Cmd) Args(name string, p *[]string) {
	if c.hasOptional {
		panic("Cmd: cannot add non-optional after optional argument")
	}
	if c.hasMulti {
		panic("Cmd: cannot have multiple repeated arguments")
	}
	c.args = append(c.args, arg{
		name:  name,
		multi: p,
	})
	c.hasMulti = true
}

// OptionalArgs defines an argument that can be present zero or more times.
func (c *Cmd) OptionalArgs(name string, p *[]string) {
	if c.hasMulti {
		panic("Cmd: cannot have multiple repeated arguments")
	}
	c.args = append(c.args, arg{
		name:     name,
		optional: true,
		multi:    p,
	})
	c.hasMulti = true
}

// Help returns a help message.
func (c *Cmd) Help() string {
	return ""
}

// Parse parses the given command-line arguments and sets values for given flags. It’s usually
// called with os.Args[1:].
func (c *Cmd) Parse(args []string) {
	err, help := c.parse(args)
	if err != nil {
		c.fail(err)
	}
	if help {
		fmt.Print(c.Help())
		os.Exit(0)
	}
}

func (c *Cmd) parse(args []string) (err error, help bool) {
	err, help, args = c.Flags.parse(args)
	if err != nil || help {
		return err, help
	}

	for _, a := range c.args {
		if len(args) == 0 {
			if !a.optional {
				return fmt.Errorf("missing %s argument", a.name), false
			} else {
				return nil, false
			}
		}
		if a.single != nil {
			*a.single = args[0]
			args = args[1:]
		} else {
			*a.multi = make([]string, len(args))
			for i, arg := range args {
				(*a.multi)[i] = arg
			}
		}
	}

	if len(args) > 0 {
		return errors.New("extra arguments on command-line"), false
	}

	return nil, false
}

func (c *Cmd) fail(err error) {
	fmt.Printf("%s: %s\n", c.name, err)
	fmt.Printf("Try '%s --help' for more information.", c.name)
	os.Exit(2)
}
