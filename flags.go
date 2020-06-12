package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

// Flags is used to define flags with and without arguments. It’s meant to be used through Cmd and
// Group.
type Flags struct {
	flags   map[string]*bool
	options map[string]option
}

type option struct {
	set func(name, value string) error
}

// Flag defines a flag without a value.
func (f *Flags) Flag(spec string, p *bool, usage string) {
	f.flags[spec] = p
}

// String defines a flag with a string value.
func (f *Flags) String(spec string, p *string, name, usage string) {
	f.options[spec] = option{
		set: func(name, value string) error {
			*p = value
			return nil
		},
	}
}

// Int defines a flag with an integer value.
func (f *Flags) Int(spec string, p *int, name, usage string) {
	f.options[spec] = option{
		set: func(name, value string) error {
			i, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid %s argument '%s'", name, value)
			}
			*p = i
			return nil
		},
	}
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
