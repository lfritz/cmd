package cmd

import (
	"reflect"
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

func ambiguousArgsTest(t *testing.T, f func(c *Cmd)) {
	defer func() {
		if recover() == nil {
			t.Errorf("Cmd did't panic for ambiguous argument sequence")
		}
	}()
	f(New("cp", func() {}))
}

func TestAmbiguousArgs(t *testing.T) {
	ambiguousArgsTest(t, func(c *Cmd) {
		c.RepeatedArg("SOURCE", new([]string))
		c.RepeatedArg("DEST", new([]string))
	})
	ambiguousArgsTest(t, func(c *Cmd) {
		c.OptionalArg("SOURCE", new(string))
		c.Arg("DEST", new(string))
		c.OptionalArg("BACKUP", new(string))
	})
	ambiguousArgsTest(t, func(c *Cmd) {
		c.RepeatedArg("SOURCE", new([]string))
		c.OptionalArg("DEST", new(string))
	})
}

func TestRepeatedArgs(t *testing.T) {
	var x, y string
	var zs []string
	command := New("command", func() {})
	command.Arg("X", &x)
	command.Arg("Y", &y)
	command.RepeatedArg("Z", &zs)
	cases := []struct {
		args                []string
		wantError, wantHelp bool
		wantX, wantY        string
		wantZs              []string
	}{
		{
			args:      []string{"x"},
			wantError: true,
		},
		{
			args:      []string{"x", "y"},
			wantError: true,
		},
		{
			args:   []string{"x", "y", "z1"},
			wantX:  "x",
			wantY:  "y",
			wantZs: []string{"z1"},
		},
		{
			args:   []string{"x", "y", "z1", "z2", "z3"},
			wantX:  "x",
			wantY:  "y",
			wantZs: []string{"z1", "z2", "z3"},
		},
		{
			args:     []string{"-h", "foo", "bar"},
			wantHelp: true,
		},
	}
	for _, c := range cases {
		x, y = "", ""
		zs = nil
		help, err := command.parse(c.args)
		if !((err != nil) == c.wantError && help == c.wantHelp) {
			errorString := "nil"
			if c.wantError {
				errorString = "non-nil"
			}
			t.Errorf("parse(%v) == %v, %v, want %v, %v",
				c.args, err, help, errorString, c.wantHelp)
		}
		if c.wantError || c.wantHelp {
			continue
		}
		if !(x == c.wantX && y == c.wantY && reflect.DeepEqual(zs, c.wantZs)) {
			t.Errorf("parse(%v) set x,y,z to %v,%v,%v, want %v,%v,%v",
				c.args, x, y, zs, c.wantX, c.wantY, c.wantZs)
		}
	}
}

func TestOptionalArgs(t *testing.T) {
	var x, y, z string
	command := New("command", func() {})
	command.OptionalArg("X", &x)
	command.OptionalArg("Y", &y)
	command.Arg("Z", &z)
	cases := []struct {
		args                []string
		wantError, wantHelp bool
		wantX, wantY, wantZ string
	}{
		{
			args:  []string{"z"},
			wantZ: "z",
		},
		{
			args:  []string{"y", "z"},
			wantY: "y",
			wantZ: "z",
		},
		{
			args:  []string{"x", "y", "z"},
			wantX: "x",
			wantY: "y",
			wantZ: "z",
		},
		{
			args:     []string{"-h", "foo", "bar"},
			wantHelp: true,
		},
		{
			args:      []string{"x", "y", "z", "a"},
			wantError: true,
		},
	}
	for _, c := range cases {
		x, y, z = "", "", ""
		help, err := command.parse(c.args)
		if !((err != nil) == c.wantError && help == c.wantHelp) {
			errorString := "nil"
			if c.wantError {
				errorString = "non-nil"
			}
			t.Errorf("parse(%v) == %v, %v, want %v, %v",
				c.args, err, help, errorString, c.wantHelp)
		}
		if c.wantError || c.wantHelp {
			continue
		}
		if !(x == c.wantX && y == c.wantY && z == c.wantZ) {
			t.Errorf("parse(%v) set x,y,z to %v,%v,%v, want %v,%v,%v",
				c.args, x, y, z, c.wantX, c.wantY, c.wantZ)
		}
	}
}
