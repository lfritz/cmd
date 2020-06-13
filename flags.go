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
	entries []*entry
}

func newFlags() Flags {
	return Flags{
		flags:   make(map[string]*bool),
		options: make(map[string]*option),
		entries: []*entry{},
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

func (e *entry) flagDefinition(maxCols int) flagDefinition {
	// add e.value where needed
	withValue := []string{}
	for _, name := range e.names {
		if e.value != "" {
			name = fmt.Sprintf("%s %s", name, e.value)
		}
		withValue = append(withValue, name)
	}

	// join lines where it makes sense
	joined := strings.Join(withValue, ", ")
	last := len(withValue) - 1
	if len([]rune(withValue[last])) > maxCols {
		return flagDefinition{
			separate: withValue,
		}
	}
	if len([]rune(joined)) > maxCols {
		return flagDefinition{
			separate: withValue[:last],
			inline:   withValue[last],
		}
	}
	return flagDefinition{
		inline: joined,
	}
}

func (e *entry) wrapUsage(maxCols int) []string {
	// split usage message into words and convert to []rune
	words := strings.Split(e.usage, " ")
	runeWords := make([][]rune, len(words))
	for i, word := range words {
		runeWords[i] = []rune(word)
	}

	// join words into lines of up to maxCols characters
	lines := []string{}
	current := []rune{}
	firstWord := true
	for _, word := range runeWords {
		if len(current)+1+len(word) > maxCols {
			lines = append(lines, string(current))
			current = word
			continue
		}
		if !firstWord {
			current = append(current, ' ')
		}
		current = append(current, word...)
		firstWord = false
	}
	lines = append(lines, string(current))

	return lines
}

func (f *Flags) printHelp(w io.Writer, columns int) {
	// set a maximum for left column size
	maxLeftCols := (columns - 4) / 2
	if maxLeftCols > 25 {
		maxLeftCols = 25
	}

	// get text for left column
	flagDefinitions := []flagDefinition{}
	for _, entry := range f.entries {
		flagDefinitions = append(flagDefinitions, entry.flagDefinition(maxLeftCols))
	}

	// find out size of left and right column
	leftCols := 0
	for _, d := range flagDefinitions {
		if d.inline == "" {
			continue
		}
		cols := len([]rune(d.inline))
		if cols > leftCols {
			leftCols = cols
		}
	}
	rightCols := columns - 4 - leftCols
	if rightCols > 50 {
		rightCols = 50
	}

	// print
	fmt.Fprintln(w, "Options:")
	for i, entry := range f.entries {
		flagDef := flagDefinitions[i]
		for _, line := range flagDef.separate {
			fmt.Fprintf(w, "  %s\n", line)
		}
		usageLines := entry.wrapUsage(rightCols)
		fmt.Fprintf(w, "  %-*s  %s\n", leftCols, flagDef.inline, usageLines[0])
		for _, line := range usageLines[1:] {
			fmt.Fprintf(w, "%*s%s\n", 2+leftCols+2, "", line)
		}
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

	f.entries = append(f.entries, &entry{
		names: names,
		usage: usage,
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

	f.entries = append(f.entries, &entry{
		names: names,
		value: name,
		usage: usage,
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
