# cmd

![Go build, lint and test](https://github.com/lfritz/cmd/workflows/Go%20build,%20lint%20and%20test/badge.svg)
[![GoDoc](https://godoc.org/github.com/lfritz/cmd?status.svg)](https://godoc.org/github.com/lfritz/cmd)
[![Go Report Card](https://goreportcard.com/badge/github.com/lfritz/cmd)](https://goreportcard.com/report/github.com/lfritz/cmd)

cmd is a simple but powerful library for building command-line programs. It supports flags,
positional arguments, and nested groups of commands.

To see how it works, let’s re-create the Unix mkdir command:

```go
package main

import (
	"os"

	"github.com/lfritz/cmd"
)

var (
	mode             string
	parents, verbose bool
	dir              []string
)

func main() {
	c := cmd.New("mkdir", run)
	c.Summary = "Create the DIRECTORY(ies), if they do not already exist."
	c.String("-m --mode", &mode, "MODE", "set file mode (as in chmod)")
	c.Flag("-p --parents", &parents, "make parent directories as needed")
	c.Flag("-v --verbose", &verbose, "print a message for each created directory")
	c.Args("DIRECTORY", &dir)
	c.Run(os.Args[1:])
}

func run() {
	// ...
}
```

If we run this program with

    go run main.go -p foo/bar

the cmd library will set `parents` to `true` and `dir` to `[]string{"foo/bar"}`, then call the `run`
function.


## Generated help

One goals of the library is to generate help message that look good regardless of the size of the
terminal. Here’s what it looks like with the example above and a narrow 50-column terminal:

```
$ go run main.go --help
Usage: mkdir [OPTION]... DIRECTORY...

Create the DIRECTORY(ies), if they do not already exist.

Options:
  -m MODE, --mode MODE  set file mode (as in chmod)
  -p, --parents         make parent directories as needed
  -v, --verbose         print a message for each created
                        directory
```

Notice that the last line has been wrapped to fit the 50-column width.


## Groups

For programs that have multiple commands (think “git status,” “git add,” etc), you create a Group
and add the commands one-by-one:

```go
package main

import (
	"os"

	"github.com/lfritz/cmd"
)

func main() {
	g := cmd.NewGroup("git")
	g.Summary = "Git is a fast, scalable, distributed revision control system."
	add := g.Command("add", add)
	add.Summary = "Add file contents to the index"
	commit := g.Command("commit", commit)
	commit.Summary = "Record changes to the repository"
	g.Run(os.Args[1:])
}

func add() {
	// ...
}

func commit() {
	// ...
}
```

Each sub-command is represented by the same type as a top-level command, so you can define flags and
arguments in the same way.
