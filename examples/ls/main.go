package main

import (
	"os"

	"github.com/lfritz/cmd"
)

var (
	long, color, all bool
	hidePattern      string
	width            int
	files            []string
)

// This example implements a subset of the “ls” command’s interface. It shows how to use flags and
// positional arguments.
func main() {
	c := cmd.New("ls")
	c.Summary = "List information about the FILEs (current directory by default)"
	c.Details = "The LS_COLORS environment variable can be used instead of --color."
	c.Flag("-l", &long, "use a long listing format")
	c.Flag("--color", &color, "colorize the output")
	c.Flag("-a --all", &all, "do not ignore entries starting with .")
	c.String("--hide", &hidePattern, "PATTERN",
		"do not list implied entries matching shell PATTERN")
	c.Int("-w --width", &width, "COLS", "set output width to COLS")
	c.OptionalArgs("FILE", &files)
	c.Parse(os.Args[1:])
	// ...
}
