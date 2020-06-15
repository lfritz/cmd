package cmd

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var helpFlags = map[string]bool{
	"-h":     true,
	"-help":  true,
	"--help": true,
}

// Flags is used to define flags with and without arguments. It’s meant to be used through Cmd and
// Group.
type Flags struct {
	// used for parsing
	flags   map[string]*bool
	options map[string]*option

	// used for help message
	defs []*definition
}

func newFlags() Flags {
	return Flags{
		flags:   make(map[string]*bool),
		options: make(map[string]*option),
		defs:    []*definition{},
	}
}

type option struct {
	set func(name, value string) error
}

type entry struct {
	names []string
	value string
	usage string
}

type flagDefinition struct {
	separate []string
	inline   string
}

func (f *Flags) usage() string {
	switch len(f.defs) {
	case 0:
		return ""
	case 1:
		return "[OPTION]"
	default:
		return "[OPTION]..."
	}
}

func (f *Flags) printDefinitions(w io.Writer, columns int) {
	if len(f.defs) > 0 {
		fmt.Fprintln(w, "Options:")
		printDefinitions(w, f.defs, columns)
		fmt.Fprintln(w)
	}
}

// Flag defines a flag without a value.
func (f *Flags) Flag(spec string, p *bool, usage string) {
	names, err := splitSpec(spec)
	if err != nil {
		panic(err.Error())
	}
	for _, name := range names {
		f.flags[name] = p
	}

	f.defs = append(f.defs, &definition{
		terms: names,
		text:  usage,
	})
}

// String defines a flag with a string value.
func (f *Flags) String(spec string, p *string, name, usage string) {
	f.addOption(spec, name, usage, func(name, value string) error {
		*p = value
		return nil
	})
}

// Int defines a flag with an integer value.
func (f *Flags) Int(spec string, p *int, name, usage string) {
	f.addOption(spec, name, usage, func(name, value string) error {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid %s argument '%s'", name, value)
		}
		*p = i
		return nil
	})
}

func (f *Flags) addOption(spec, name, usage string, set func(name, value string) error) {
	names, err := splitSpec(spec)
	if err != nil {
		panic(err.Error())
	}
	op := &option{
		set: set,
	}
	for _, name := range names {
		f.options[name] = op
	}

	// add e.value where needed
	terms := []string{}
	for _, n := range names {
		terms = append(terms, fmt.Sprintf("%s %s", n, name))
	}

	f.defs = append(f.defs, &definition{
		terms: terms,
		text:  usage,
	})
}

var splitRe = regexp.MustCompile(`^--?[^-]`)

func splitSpec(spec string) ([]string, error) {
	fail := func() ([]string, error) {
		return nil, fmt.Errorf("invalid spec: %s", spec)
	}
	parts := strings.Split(spec, " ")
	if len(parts) == 0 {
		return fail()
	}
	for _, p := range parts {
		if !splitRe.MatchString(p) {
			return fail()
		}
	}
	return parts, nil
}

func (f *Flags) parse(args []string) (err error, help bool, following []string) {
	for len(args) > 0 {
		a := args[0]
		if !isFlag(a) {
			break
		}
		args = args[1:]

		a, value := splitFlag(a)
		if value != "" {
			_, ok := f.flags[a]
			if ok {
				return fmt.Errorf("%s does not take a value", a), false, nil
			}
			o, ok := f.options[a]
			if !ok {
				return fmt.Errorf("unrecognized flag %s", a), false, nil
			}
			err := o.set(a, value)
			if err != nil {
				return err, false, nil
			}
			continue
		}

		if helpFlags[a] {
			return nil, true, nil
		}

		ptr, ok := f.flags[a]
		if ok {
			*ptr = true
			continue
		}

		o, ok := f.options[a]
		if ok {
			if len(args) == 0 {
				return fmt.Errorf("missing value for argument %s", a), false, nil
			}
			err := o.set(a, args[0])
			args = args[1:]
			if err != nil {
				return err, false, nil
			}
			continue
		}

		return fmt.Errorf("unrecognized flag %s", a), false, nil
	}

	return nil, false, args
}

func isFlag(s string) bool {
	if s == "" {
		return false
	}
	return s[0] == '-'
}

func splitFlag(s string) (string, string) {
	slice := strings.SplitN(s, "=", 2)
	if len(slice) == 1 {
		return slice[0], ""
	}
	return slice[0], slice[1]
}
