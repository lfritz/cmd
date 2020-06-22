package cmd

import (
	"fmt"
	"regexp"
	"strings"
)

type flagParser struct {
	allowPositional bool
	flags           map[string]*bool
	options         map[string]*option
}

type option struct {
	set func(name, value string) error
}

func newFlagParser(allowPositional bool) *flagParser {
	return &flagParser{
		allowPositional: allowPositional,
		flags:           make(map[string]*bool),
		options:         make(map[string]*option),
	}
}

func (p *flagParser) addFlag(spec string, ptr *bool) {
	names, err := splitSpec(spec)
	if err != nil {
		panic(err.Error())
	}
	for _, name := range names {
		p.flags[name] = ptr
	}
}

func (p *flagParser) addOption(spec string, set func(name, value string) error) {
	names, err := splitSpec(spec)
	if err != nil {
		panic(err.Error())
	}
	op := &option{
		set: set,
	}
	for _, name := range names {
		p.options[name] = op
	}
}

var singleDashMultiCharRe = regexp.MustCompile(`^-[^-].`)

func (p *flagParser) parse(args []string, singleDashMode, allowVersion bool) (remaining []string, help, version bool, err error) {
	helpFlags := make(map[string]bool)
	versionFlags := make(map[string]bool)
	if singleDashMode {
		helpFlags["-h"] = true
		helpFlags["--help"] = true
		versionFlags["-v"] = true
		versionFlags["--version"] = true
	} else {
		helpFlags["-help"] = true
		versionFlags["-version"] = true
	}

	remaining = []string{}
	for len(args) != 0 {
		a := args[0]
		if !isFlag(a) && !p.allowPositional {
			// a must be a sub-command
			remaining = args
			return
		}
		args = args[1:]
		if !isFlag(a) {
			// a must be a positional argument
			remaining = append(remaining, a)
			continue
		}
		if a == "---" {
			// "---" means "treat the rest of the command-line as positional arguments"
			remaining = args
			return
		}
		if name, value := splitFlag(a); value != "" {
			// a is of the form "--name=value"
			_, ok := p.flags[a]
			if ok {
				err = fmt.Errorf("%s does not take a value", a)
				return
			}
			o, ok := p.options[name]
			if !ok {
				err = fmt.Errorf("unrecognized flag %s", name)
				return
			}
			err = o.set(a, value)
			if err != nil {
				return
			}
			continue
		}
		if !singleDashMode && singleDashMultiCharRe.MatchString(a) {
			// a is multiple single-character flags (think "ls -la")
			split := []string{}
			for _, c := range a[1:] {
				split = append(split, fmt.Sprintf("-%c", c))
			}
			args = append(split, args...)
			continue
		}
		if ptr, ok := p.flags[a]; ok {
			// a is a flag
			*ptr = true
			continue
		}
		if o, ok := p.options[a]; ok {
			// a is an option
			if len(args) == 0 {
				err = fmt.Errorf("missing value for argument %s", a)
				return
			}
			err = o.set(a, args[0])
			if err != nil {
				return
			}
			args = args[1:]
			continue
		}
		if helpFlags[a] {
			// a is --help or equivalent
			help = true
			return
		}
		if allowVersion && versionFlags[a] {
			// a is --version or equivalent
			version = true
			return
		}
		err = fmt.Errorf("unrecognized flag '%s'", a)
		return
	}

	return
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
