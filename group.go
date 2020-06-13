package cmd

import (
	"fmt"
	"os"
)

// A Group represents a group of commands. Groups can be nested arbitrarily.
type Group struct {
	Flags
	Summary, Details string
	name             string
	groups           map[string]*Group
	commands         map[string]*Cmd
}

// NewGroup returns a new group of commands with the specified name.
func NewGroup(name string) *Group {
	return nil
}

// Command adds a command. The given function will be called if this command is selected.
func (g *Group) Command(name string, f func()) *Cmd {
	command := New(name, f)
	g.commands[name] = command
	return command
}

// Group adds a sub-group.
func (g *Group) Group(name string) *Group {
	group := NewGroup(name)
	g.groups[name] = group
	return group
}

// Help returns a help message.
func (g *Group) Help() string {
	return ""
}

// Run parses the given command-line arguments, sets values for given flags and calls the function
// for the selected command. It’s usually called with os.Args[1:].
func (g *Group) Run(args []string) {
	// call Flags.parse
	err, help, args := g.Flags.parse(args)
	if err != nil {
		g.fail(err.Error())
	}
	if help {
		g.help()
	}

	// select group or command
	if len(args) == 0 {
		g.fail("command expected")
	}
	a := args[0]
	args = args[1:]
	if group, ok := g.groups[a]; ok {
		group.Run(args)
		return
	}
	if command, ok := g.commands[a]; ok {
		command.Run(args)
		return
	}
	g.fail(fmt.Sprintf("'%s' is not a %s command", a, g.name))
}

func (g *Group) fail(msg string) {
	fmt.Printf("%s: %s\n", g.name, msg)
	fmt.Printf("Try '%s help' for more information.", g.name)
	os.Exit(2)
}

func (g *Group) help() {
	fmt.Print(g.Help())
	os.Exit(0)
}
