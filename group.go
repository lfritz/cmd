package cmd

import (
	"fmt"
	"os"
	"strings"
)

// A Group represents a group of commands. Groups can be nested arbitrarily.
//
// The Summary and Details fields are printed at the beginning and end, respectively, of the help
// message. They won’t be printed if left empty.
type Group struct {
	Flags
	Summary, Details string
	name             string
	groups           map[string]*Group
	commands         map[string]*Cmd
}

// NewGroup returns a new group of commands with the specified name.
func NewGroup(name string) *Group {
	return &Group{
		Flags:    newFlags(),
		name:     name,
		groups:   make(map[string]*Group),
		commands: make(map[string]*Cmd),
	}
}

// Command adds a command.
func (g *Group) Command(name string, f func()) *Cmd {
	command := New(fmt.Sprintf("%s %s", g.name, name), f)
	g.commands[name] = command
	return command
}

// Group adds a sub-group.
func (g *Group) Group(name string) *Group {
	group := NewGroup(fmt.Sprintf("%s %s", g.name, name))
	g.groups[name] = group
	return group
}

func (g *Group) errorAndExit(msg string) {
	w := os.Stderr
	fmt.Fprintf(w, "%s: %s\n", g.name, msg)
	fmt.Fprintf(w, "Try '%s help' for more information.\n", g.name)
	os.Exit(2)
}

func (g *Group) helpAndExit() {
	fmt.Fprintf(os.Stdout, g.Help())
	os.Exit(0)
}

// Help returns a help message.
func (g *Group) Help() string {
	defs := []*definitionList{
		{
			title:       "Options",
			definitions: g.Flags.defs,
		},
		{
			title:       "Groups",
			definitions: g.groupDefinitions(),
		},
		{
			title:       "Commands",
			definitions: g.commandDefinitions(),
		},
	}
	return formatHelp(g.usage(), g.Summary, g.Details, defs)
}

func (g *Group) summary() string {
	return g.Summary
}

func (g *Group) details() string {
	return g.Details
}

func (g *Group) groupDefinitions() []*definition {
	defs := []*definition{}
	for name, g := range g.groups {
		defs = append(defs, &definition{
			terms: []string{name},
			text:  g.Summary,
		})
	}
	return defs
}

func (g *Group) commandDefinitions() []*definition {
	defs := []*definition{}
	for name, c := range g.commands {
		defs = append(defs, &definition{
			terms: []string{name},
			text:  c.Summary,
		})
	}
	return defs
}

func (g *Group) usage() string {
	line := []string{"Usage:", g.name}
	if s := g.Flags.usage(); s != "" {
		line = append(line, s)
	}
	groupOrCommand := []string{}
	if len(g.groups) > 0 {
		groupOrCommand = append(groupOrCommand, "GROUP")
	}
	if len(g.commands) > 0 {
		groupOrCommand = append(groupOrCommand, "COMMAND")
	}
	line = append(line, strings.Join(groupOrCommand, " | "))
	return strings.Join(line, " ")
}

// Run parses the given command-line arguments, sets values for given flags and calls the function
// for the selected command. It’s usually called with os.Args[1:].
func (g *Group) Run(args []string) {
	g.run(args, false)
}

func (g *Group) run(args []string, helpMode bool) {
	// call Flags.parse
	help, args, err := g.Flags.parse(args)
	if err != nil {
		g.errorAndExit(err.Error())
	}
	if help {
		g.helpAndExit()
	}

	// select group or command
	if len(args) == 0 {
		if helpMode {
			g.helpAndExit()
		}
		g.errorAndExit("command expected")
	}
	a, args := args[0], args[1:]
	if a == "help" {
		g.run(args, true)
		return
	}
	if group, ok := g.groups[a]; ok {
		group.run(args, helpMode)
		return
	}
	if command, ok := g.commands[a]; ok {
		if helpMode {
			command.helpAndExit()
		} else {
			command.Run(args)
		}
		return
	}
	g.errorAndExit(fmt.Sprintf("'%s' is not a %s command", a, g.name))
}
