package cmd

import "fmt"

// This example implements a subset of the Unix ls command’s interface.
func ExampleCmd() {
	var (
		long, color, all bool
		hidePattern      string
		files            []string
	)

	c := New("ls", func() {})
	c.Summary = "List information about the FILEs (current directory by default)."
	c.Details = "The LS_COLORS environment variable can be used instead of --color."
	c.Flag("-l", &long, "use a long listing format")
	c.Flag("--color", &color, "colorize the output")
	c.Flag("-a --all", &all, "do not ignore entries starting with .")
	c.String("--hide", &hidePattern, "PATTERN",
		"do not list implied entries matching shell PATTERN")
	c.OptionalArgs("FILE", &files)
	fmt.Print(c.Help())
	// Output:
	// Usage: ls [OPTION]... [FILE]...
	//
	// List information about the FILEs (current directory by default).
	//
	// Options:
	//   -l              use a long listing format
	//   --color         colorize the output
	//   -a, --all       do not ignore entries starting with .
	//   --hide PATTERN  do not list implied entries matching shell PATTERN
	//
	// The LS_COLORS environment variable can be used instead of --color.
}

// This implements a subset of the git command’s interface.
func ExampleGroup() {
	var (
		path   string
		status struct {
			short    bool
			pathspec []string
		}
		init struct {
			quiet, bare bool
		}
	)

	git := NewGroup("git")
	git.Summary = "Git is a fast, scalable, distributed revision control system."
	git.String("-C", &path, "PATH", "Run as if git was started in PATH.")

	gitStatus := git.Command("status", func() {})
	gitStatus.Summary = "Show the working tree status"
	gitStatus.Flag("-s --short", &status.short, "Give the output in the short-format.")
	gitStatus.OptionalArgs("PATHSPEC...", &status.pathspec)

	gitInit := git.Command("init", func() {})
	gitInit.Summary = "Create an empty Git repository"
	gitInit.Flag("-q --quiet", &init.quiet, "Only print error and warning messages.")
	gitInit.Flag("--bare", &init.bare, "Create a bare repository.")

	gitRemote := git.Group("remote")
	gitRemote.Summary = "Manage the set of repositories you track"
	gitRemote.Command("add", func() {})
	gitRemote.Command("rename", func() {})
	gitRemote.Command("remove", func() {})

	fmt.Print(git.Help())
	// Output:
	// Usage: git [OPTION] GROUP | COMMAND
	//
	// Git is a fast, scalable, distributed revision control system.
	//
	// Options:
	//   -C PATH  Run as if git was started in PATH.
	//
	// Groups:
	//   remote  Manage the set of repositories you track
	//
	// Commands:
	//   status  Show the working tree status
	//   init    Create an empty Git repository
}
