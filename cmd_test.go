package cmd

import (
	"testing"
)

func TestCmdUsage(t *testing.T) {
	c := New("cp", func() {})
	c.Flag("-f", new(bool), "")
	c.Flag("-r", new(bool), "")
	c.RepeatedArg("SOURCE", new([]string))
	c.Arg("DEST", new(string))
	got := c.usage()
	want := "Usage: cp [OPTION]... SOURCE... DEST"
	if got != want {
		t.Errorf("usage == `%s`, want `%s`", got, want)
	}
}
