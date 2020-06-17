package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var helpFlags = map[string]bool{
	"-h":     true,
	"-help":  true,
	"--help": true,
}

// Flags is used to define flags with and without arguments. It’s embedded in Cmd and Group; you
// usually call its methods directly on those types.
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

// Float defines a flag with a float64 value. See strconv.ParseFloat for the format it recognizes.
func (f *Flags) Float(spec string, p *float64, name, usage string) {
	f.addOption(spec, name, usage, func(name, value string) error {
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
func (f *Flags) Duration(spec string, p *time.Duration, name, usage string) {
	f.addOption(spec, name, usage, func(name, value string) error {
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
func (f *Flags) Metric(spec string, p *int, name, usage string) {
	f.addOption(spec, name, usage, func(name, value string) error {
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
func (f *Flags) Bytes(spec string, p *int, name, usage string) {
	f.addOption(spec, name, usage, func(name, value string) error {
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
