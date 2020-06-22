package cmd

import (
	"errors"
	"fmt"
)

// argSequence is a simplified view of the positional arguments added so far, using the constants
// defined below.
type argSequence int

const (
	argSequenceInitial         = iota // no positional args so far
	argSequenceRegular                // one or more regular args (not optional or repeated)
	argSequenceRegularOptional        // one or more regular, then one or more optional
	argSequenceRegularRepeated        // one or more regular, then one repeated
	argSequenceOptional               // one or more optional
	argSequenceRegularAtEnd           // one or more optional OR one repeated, then regular args
)

type argParser struct {
	state argSequence
	args  []positionalArgument
}

type positionalArgument struct {
	name     string
	optional bool
	single   *string
	slice    *[]string
}

func (a positionalArgument) isOptional() bool {
	return a.optional
}
func (a positionalArgument) isRepeated() bool {
	return a.single == nil
}

func newArgParser() *argParser {
	return &argParser{
		args: []positionalArgument{},
	}
}

func (p *argParser) add(name string, ptr *string) {
	switch p.state {
	case argSequenceInitial, argSequenceRegular:
		p.addArg(name, false, ptr, nil)
		p.state = argSequenceRegular
	case argSequenceOptional, argSequenceRegularAtEnd:
		p.addArg(name, false, ptr, nil)
		p.state = argSequenceRegularAtEnd
	default:
		ambiguousArgs()
	}
}

func (p *argParser) addOptional(name string, ptr *string) {
	switch p.state {
	case argSequenceInitial, argSequenceOptional:
		p.addArg(name, true, ptr, nil)
		p.state = argSequenceOptional
	case argSequenceRegular, argSequenceRegularOptional:
		p.addArg(name, true, ptr, nil)
		p.state = argSequenceRegularOptional
	default:
		ambiguousArgs()
	}
}

func (p *argParser) addRepeated(name string, ptr *[]string, optional bool) {
	switch p.state {
	case argSequenceInitial:
		p.addArg(name, optional, nil, ptr)
		p.state = argSequenceRegularAtEnd
	case argSequenceRegular:
		p.addArg(name, optional, nil, ptr)
		p.state = argSequenceRegularRepeated
	default:
		ambiguousArgs()
	}
}

func (p *argParser) addArg(name string, optional bool, single *string, slice *[]string) {
	p.args = append(p.args, positionalArgument{
		name:     name,
		optional: optional,
		single:   single,
		slice:    slice,
	})
}

func (p *argParser) parse(args []string) error {
	if p.state == argSequenceRegularAtEnd {
		return p.parseBackward(args)
	}

	for _, pa := range p.args {
		if len(args) == 0 && pa.isOptional() {
			return nil
		}
		if len(args) == 0 {
			return fmt.Errorf("missing %s argument", pa.name)
		}
		if pa.isRepeated() {
			*pa.slice = args
			return nil
		}
		*pa.single = args[0]
		args = args[1:]
	}

	if len(args) != 0 {
		return errors.New("extra arguments on command-line")
	}

	return nil
}

func (p *argParser) parseBackward(args []string) error {
	for i := len(p.args) - 1; i >= 0; i-- {
		pa := p.args[i]
		if len(args) == 0 && pa.isOptional() {
			return nil
		}
		if len(args) == 0 {
			return fmt.Errorf("missing %s argument", pa.name)
		}
		if pa.isRepeated() {
			*pa.slice = args
			return nil
		}
		*pa.single = args[len(args)-1]
		args = args[:len(args)-1]
	}

	if len(args) != 0 {
		return errors.New("extra arguments on command-line")
	}

	return nil
}

func ambiguousArgs() {
	panic("Cmd: ambiguous sequence of positional arguments")
}
