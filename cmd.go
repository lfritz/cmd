// Package cmd implements a command-line parser.
package cmd

// A Cmd represents a command with command-line flags and arguments.
type Cmd struct {
	Flags
	Summary, Details string
}

// New returns a new command with the specified name.
func New(name string) *Cmd {
	return nil
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
