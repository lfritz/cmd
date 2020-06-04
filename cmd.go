// Package cmd implements a command-line parser.
package cmd

// A Cmd represents a set of command-line flags and arguments.
type Cmd struct{}

// New returns a new command with the specified name. If non-empty, the summary is printed before
// and the details after the flags in the help message.
func New(name, summary, details string) *Cmd {
	return nil
}

// Flag defines a flag without a value.
func (c *Cmd) Flag(spec string, p *bool, usage string) {
}

// String defines a flag with a string value.
func (c *Cmd) String(spec string, p *string, name, usage string) {
}

// Int defines a flag with an int value.
func (c *Cmd) Int(spec string, p *int, name, usage string) {
}

// Arg defines a positional argument.
func (c *Cmd) Arg(name string, p *string) {
}

// OptionalArg defines an optional positional argument.
func (c *Cmd) OptionalArg(name string, p *string) {
}

// Args defines an argument that can be present one or more times.
func (c *Cmd) Args(name string, p *[]string) {
}

// OptionalArgs defines an argument that can be present zero or more times.
func (c *Cmd) OptionalArgs(name string, p *[]string) {
}

// Help returns a help message.
func (c *Cmd) Help() string {
	return ""
}

// Parse parses the given command-line arguments and sets values for given flags. It’s usually
// called with with os.Args[1:].
func (c *Cmd) Parse(args []string) {
}
