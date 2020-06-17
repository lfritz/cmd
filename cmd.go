// Package cmd implements a command-line parser.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// A Cmd represents a command with flags and positional arguments.
//
// The *Arg methods will panic if adding the argument would mean the command-line can’t be parsed
// unambiguously anymore. For example, successive calls to OptionalArg, Arg, OptionalArg mean that
// if the use passes two arguments, there’s no way to tell which optional argument they’re trying to
// specify.
//
// The Summary and Details fields are printed at the beginning and end, respectively, of the help
// message. They won’t be printed if left empty.
type Cmd struct {
	Flags
	Summary, Details string
	name             string
	f                func()
	args             []arg
	argsState        int
}

type arg struct {
	name     string
	optional bool
	single   *string
	multi    *[]string
}

const (
	argsInitial = iota
	argsRegular
	argsRegularOptional
	argsRegularMulti
	argsMulti
	argsMultiRegular
	argsOptinal
	argsOptinalRegular
)

// New returns a new command that calls the given function after parsing arguments. The name is used
// in help and error messages.
func New(name string, f func()) *Cmd {
	return &Cmd{
		Flags: newFlags(),
		name:  name,
		f:     f,
	}
}

func ambiguousArgs() {
	panic("Cmd: ambiguous sequence of positional arguments")
}

// Arg defines a positional argument.
func (c *Cmd) Arg(name string, p *string) {
	switch c.argsState {
	case argsInitial, argsRegular:
		c.argsState = argsRegular
	case argsMulti:
		c.argsState = argsMultiRegular
	case argsOptinal:
		c.argsState = argsOptinalRegular
	default:
		ambiguousArgs()
	}
	c.args = append(c.args, arg{
		name:   name,
		single: p,
	})
}

// OptionalArg defines an optional positional argument.
func (c *Cmd) OptionalArg(name string, p *string) {
	switch c.argsState {
	case argsInitial, argsOptinal:
		c.argsState = argsOptinal
	case argsRegular:
		c.argsState = argsRegularOptional
	default:
		ambiguousArgs()
	}
	c.args = append(c.args, arg{
		name:     name,
		optional: true,
		single:   p,
	})
}

// Args defines an argument that can be present one or more times.
func (c *Cmd) Args(name string, p *[]string) {
	c.addArgs(name, p, false)
}

// OptionalArgs defines an argument that can be present zero or more times.
func (c *Cmd) OptionalArgs(name string, p *[]string) {
	c.addArgs(name, p, true)
}

func (c *Cmd) addArgs(name string, p *[]string, optional bool) {
	switch c.argsState {
	case argsInitial:
		c.argsState = argsMulti
	case argsRegular:
		c.argsState = argsRegularMulti
	default:
		ambiguousArgs()
	}
	c.args = append(c.args, arg{
		name:     name,
		optional: optional,
		multi:    p,
	})
}

func (c *Cmd) errorAndExit(err error) {
	w := os.Stderr
	fmt.Fprintf(w, "%s: %s\n", c.name, err)
	fmt.Fprintf(w, "Try '%s --help' for more information.", c.name)
	os.Exit(2)
}

func (c *Cmd) helpAndExit() {
	fmt.Fprintf(os.Stdout, c.Help())
	os.Exit(0)
}

// Help returns a help message.
func (c *Cmd) Help() string {
	defs := []*definitionList{
		{
			title:       "Options",
			definitions: c.Flags.defs,
		},
	}
	return formatHelp(c.usage(), c.Summary, c.Details, defs)
}

func (c *Cmd) usage() string {
	line := []string{"Usage:", c.name}
	if s := c.Flags.usage(); s != "" {
		line = append(line, s)
	}
	for _, arg := range c.args {
		var s string
		if arg.optional {
			s = fmt.Sprintf("[%s]", arg.name)
		} else {
			s = fmt.Sprintf("%s", arg.name)
		}
		if arg.multi != nil {
			s += "..."
		}
		line = append(line, s)
	}
	return strings.Join(line, " ")
}

// Run parses the given command-line arguments, sets values for given flags and runs the function
// provided to New. It’s usually called with os.Args[1:].
func (c *Cmd) Run(args []string) {
	err, help := c.parse(args)
	if err != nil {
		c.errorAndExit(err)
	}
	if help {
		c.helpAndExit()
	}
	c.f()
}

func (c *Cmd) parse(args []string) (err error, help bool) {
	// parse flags
	err, help, args = c.Flags.parse(args)
	if err != nil || help {
		return err, help
	}

	if c.argsState >= argsMulti {
		// parse positional arguments in reverse order
		for i := len(c.args) - 1; i >= 0; i-- {
			a := c.args[i]
			if len(args) == 0 {
				if !a.optional {
					return fmt.Errorf("missing %s argument", a.name), false
				} else {
					return nil, false
				}
			}
			if a.single != nil {
				*a.single = args[len(args)-1]
				args = args[:len(args)-1]
			} else {
				*a.multi = make([]string, len(args))
				for i, arg := range args {
					(*a.multi)[i] = arg
				}
				args = nil
			}
		}
	} else {
		// parse positional arguments in-order
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
				args = nil
			}
		}
	}

	if len(args) > 0 {
		return errors.New("extra arguments on command-line"), false
	}

	return nil, false
}
