package cmd

// A Group represents a group of commands. Groups can be nested arbitrarily.
type Group struct {
	Flags
	Summary, Details string
}

// NewGroup returns a new group of commands with the specified name.
func NewGroup(name string) *Group {
	return nil
}

// Command adds a command. The given function will be called if this command is selected.
func (g *Group) Command(name string, f func()) *Cmd {
	return nil
}

// Group adds a sub-group.
func (g *Group) Group(name string) *Group {
	return nil
}

// Help returns a help message.
func (g *Group) Help() string {
	return ""
}

// Parse parses the given command-line arguments, sets values for given flags and calls the function
// for the selected command. It’s usually called with with os.Args[1:].
func (g *Group) Run(args []string) {
}
