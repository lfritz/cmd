package cmd

// A Group represents a group of commands. Groups can be nested arbitrarily.
type Group struct {
	Cmd
}

// NewGroup returns a new group of commands with the specified name. If non-empty, the summary is
// printed before and the details after the flags in the help message.
func NewGroup(name, summary, details string) *Group {
	return nil
}

// Command adds a command. The given function will be called if this command is selected.
func (g *Group) Command(command *Cmd, f func()) {
}

// Group adds a sub-group.
func (g *Group) Group(group *Group) {
}

// Help returns a help message.
func (g *Group) Help() string {
	return ""
}

// Parse parses the given command-line arguments, sets values for given flags and calls the function
// for the selected command. It’s usually called with with os.Args[1:].
func (g *Group) Run() {
}
