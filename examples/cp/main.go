package main

import (
	"fmt"
	"os"

	"github.com/lfritz/cmd"
)

var (
	source []string
	dest   string
)

// This example implements a subset of the “cp” comnand’s interface.
func main() {
	c := cmd.New("cp", run)
	c.Args("SOURCE", &source)
	c.Arg("DEST", &dest)
	c.Run(os.Args[1:])
}

func run() {
	fmt.Println("source:", source)
	fmt.Println("dest:  ", dest)
}
