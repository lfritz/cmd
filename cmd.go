// Package cmd implements a command-line parser.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	flagParser                *flagParser
	Summary, Details, Version string
	name                      string
	f                         func()
	args                      []positionalArgument
	argsState                 int

	// used for help message
	optionDefinitions []*definition
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
		flagParser:        newFlagParser(true),
		name:              name,
		f:                 f,
		optionDefinitions: []*definition{},
	}
}

// Flag defines a flag without a value.
func (c *Cmd) Flag(spec string, p *bool, usage string) {
	names, err := splitSpec(spec)
	if err != nil {
		panic(err.Error())
	}
	c.flagParser.addFlag(names, p)
	c.optionDefinitions = append(c.optionDefinitions, &definition{
		terms: names,
		text:  usage,
	})
}

// String defines a flag with a string value.
func (c *Cmd) String(spec string, p *string, name, usage string) {
	c.addOption(spec, name, usage, func(name, value string) error {
		*p = value
		return nil
	})
}

// Int defines a flag with an integer value.
func (c *Cmd) Int(spec string, p *int, name, usage string) {
	c.addOption(spec, name, usage, func(name, value string) error {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid %s argument '%s'", name, value)
		}
		*p = i
		return nil
	})
}

// Float defines a flag with a float64 value. See strconv.ParseFloat for the format it recognizes.
func (c *Cmd) Float(spec string, p *float64, name, usage string) {
	c.addOption(spec, name, usage, func(name, value string) error {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid %s argument '%s'", name, value)
		}
		*p = f
		return nil
	})
}

// Duration defines a flag with a time.Duration value. See time.ParseDuration for the format it
// recognizes.
func (c *Cmd) Duration(spec string, p *time.Duration, name, usage string) {
	c.addOption(spec, name, usage, func(name, value string) error {
		d, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid %s argument '%s'", name, value)
		}
		*p = d
		return nil
	})
}

// Metric defines a flag with an integer value that allows the user to use metric suffixes, for
// example “5k“ for 5000. Both lower-case and upper-case suffixes work.
func (c *Cmd) Metric(spec string, p *int, name, usage string) {
	c.addOption(spec, name, usage, func(name, value string) error {
		i, ok := parseWithSuffix(value, metricSuffixMap)
		if !ok {
			return fmt.Errorf("invalid %s argument '%s'", name, value)
		}
		*p = i
		return nil
	})
}

var metricSuffixMap = map[string]int{
	"k": 1000,
	"m": 1000000,
	"g": 1000000000,
	"t": 1000000000000,
	"p": 1000000000000000,
	"e": 1000000000000000000,
}

// Bytes defines a flag with an integer value that allows the user to use binary suffixes, for
// example “5k“ for 5*1024. Both lower-case and upper-case suffixes work.
func (c *Cmd) Bytes(spec string, p *int, name, usage string) {
	c.addOption(spec, name, usage, func(name, value string) error {
		i, ok := parseWithSuffix(value, bytesSuffixMap)
		if !ok {
			return fmt.Errorf("invalid %s argument '%s'", name, value)
		}
		*p = i
		return nil
	})
}

var bytesSuffixMap = map[string]int{
	"k": 1 << 10,
	"m": 1 << 20,
	"g": 1 << 30,
	"t": 1 << 40,
	"p": 1 << 50,
	"e": 1 << 60,
}

var suffixRe = regexp.MustCompile(`^(-?[0-9]+)([a-zA-Z])?$`)

func parseWithSuffix(s string, suffixMap map[string]int) (i int, ok bool) {
	match := suffixRe.FindStringSubmatch(s)
	if match == nil {
		return 0, false
	}
	number, suffix := match[1], match[2]

	i, err := strconv.Atoi(number)
	if err != nil {
		return 0, false
	}

	factor, ok := suffixMap[strings.ToLower(suffix)]
	if !ok {
		return 0, false
	}

	return i * factor, true
}

func (c *Cmd) addOption(spec, name, usage string, set func(name, value string) error) {
	names, err := splitSpec(spec)
	if err != nil {
		panic(err.Error())
	}
	c.flagParser.addOption(names, set)

	// add e.value where needed
	terms := []string{}
	for _, n := range names {
		terms = append(terms, fmt.Sprintf("%s %s", n, name))
	}

	c.optionDefinitions = append(c.optionDefinitions, &definition{
		terms: terms,
		text:  usage,
	})
}

var splitRe = regexp.MustCompile(`^--?[^-]`)

func splitSpec(spec string) ([]string, error) {
	errorMessage := fmt.Errorf("invalid spec: %s", spec)
	parts := strings.Split(spec, " ")
	if len(parts) == 0 {
		return nil, errorMessage
	}
	for _, p := range parts {
		if !splitRe.MatchString(p) {
			return nil, errorMessage
		}
	}
	return parts, nil
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
	c.args = append(c.args, positionalArgument{
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
	c.args = append(c.args, positionalArgument{
		name:     name,
		optional: true,
		single:   p,
	})
}

// RepeatedArg defines an argument that can be present one or more times.
func (c *Cmd) RepeatedArg(name string, p *[]string) {
	c.addArgs(name, p, false)
}

// OptionalRepeatedArg defines an argument that can be present zero or more times.
func (c *Cmd) OptionalRepeatedArg(name string, p *[]string) {
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
	c.args = append(c.args, positionalArgument{
		name:     name,
		optional: optional,
		slice:    p,
	})
}

func (c *Cmd) errorAndExit(err error) {
	w := os.Stderr
	fmt.Fprintf(w, "%s: %s\n", c.name, err)
	fmt.Fprintf(w, "Try '%s --help' for more information.\n", c.name)
	os.Exit(2)
}

func (c *Cmd) helpAndExit() {
	fmt.Fprintf(os.Stdout, c.Help())
	os.Exit(0)
}

func (c *Cmd) versionAndExit() {
	fmt.Fprintf(os.Stdout, c.formatVersion())
	os.Exit(0)
}

// Help returns a help message.
func (c *Cmd) Help() string {
	defs := []*definitionList{
		{
			title:       "Options",
			definitions: c.optionDefinitions,
		},
	}
	return formatHelp(c.usage(), c.Summary, c.Details, defs)
}

func (c *Cmd) formatVersion() string {
	return fmt.Sprintf("%s %s", c.name, c.Version)
}

func (c *Cmd) usage() string {
	line := []string{"Usage:", c.name}
	if s := c.flagsUsage(); s != "" {
		line = append(line, s)
	}
	for _, arg := range c.args {
		var s string
		if arg.optional {
			s = fmt.Sprintf("[%s]", arg.name)
		} else {
			s = fmt.Sprintf("%s", arg.name)
		}
		if arg.slice != nil {
			s += "..."
		}
		line = append(line, s)
	}
	return strings.Join(line, " ")
}

func (c *Cmd) flagsUsage() string {
	switch len(c.optionDefinitions) {
	case 0:
		return ""
	case 1:
		return "[OPTION]"
	default:
		return "[OPTION]..."
	}
}

// Run parses the given command-line arguments, sets values for given flags and runs the function
// provided to New. It’s usually called with os.Args[1:].
func (c *Cmd) Run(args []string) {
	help, version, err := c.parse(args)
	if err != nil {
		c.errorAndExit(err)
	}
	if help {
		c.helpAndExit()
	}
	if version {
		c.versionAndExit()
	}
	c.f()
}

func (c *Cmd) parse(args []string) (help, version bool, err error) {
	// parse flags
	allowVersion := c.Version != ""
	args, help, version, err = c.flagParser.parse(args, allowVersion)
	if err != nil || help || version {
		return help, version, err
	}

	if c.argsState >= argsMulti {
		// parse positional arguments in reverse order
		for i := len(c.args) - 1; i >= 0; i-- {
			a := c.args[i]
			if len(args) == 0 {
				if !a.optional {
					return false, false, fmt.Errorf("missing %s argument", a.name)
				}
				return false, false, nil
			}
			if a.single != nil {
				*a.single = args[len(args)-1]
				args = args[:len(args)-1]
			} else {
				*a.slice = make([]string, len(args))
				for i, arg := range args {
					(*a.slice)[i] = arg
				}
				args = nil
			}
		}
	} else {
		// parse positional arguments in-order
		for _, a := range c.args {
			if len(args) == 0 {
				if !a.optional {
					return false, false, fmt.Errorf("missing %s argument", a.name)
				}
				return false, false, nil
			}
			if a.single != nil {
				*a.single = args[0]
				args = args[1:]
			} else {
				*a.slice = make([]string, len(args))
				for i, arg := range args {
					(*a.slice)[i] = arg
				}
				args = nil
			}
		}
	}

	if len(args) > 0 {
		return false, false, errors.New("extra arguments on command-line")
	}

	return false, false, nil
}
