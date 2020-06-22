package cmd

import (
	"reflect"
	"testing"
)

func ambiguousArgsTest(t *testing.T, f func(p *argParser)) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Helper()
			t.Errorf("argParser did't panic for ambiguous argument sequence")
		}
	}()
	f(newArgParser())
}

func TestAmbiguousArgs(t *testing.T) {
	ambiguousArgsTest(t, func(p *argParser) {
		p.addRepeated("SOURCE", new([]string), false)
		p.addRepeated("DEST", new([]string), false)
	})
	ambiguousArgsTest(t, func(p *argParser) {
		p.addOptional("SOURCE", new(string))
		p.add("DEST", new(string))
		p.addOptional("BACKUP", new(string))
	})
	ambiguousArgsTest(t, func(p *argParser) {
		p.addRepeated("SOURCE", new([]string), false)
		p.addOptional("DEST", new(string))
	})
}

func TestRepeatedArgs(t *testing.T) {
	var x, y string
	var zs []string
	p := newArgParser()
	p.add("X", &x)
	p.add("Y", &y)
	p.addRepeated("Z", &zs, false)
	cases := []struct {
		args         []string
		wantError    bool
		wantX, wantY string
		wantZs       []string
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
	}
	for _, c := range cases {
		x, y = "", ""
		zs = nil
		err := p.parse(c.args)
		if (err != nil) != c.wantError {
			errorString := "nil"
			if c.wantError {
				errorString = "non-nil"
			}
			t.Errorf("parse(%v) == %v, want %v", c.args, err, errorString)
		}
		if c.wantError {
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
	p := newArgParser()
	p.addOptional("X", &x)
	p.addOptional("Y", &y)
	p.add("Z", &z)
	cases := []struct {
		args                []string
		wantError           bool
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
			args:      []string{"x", "y", "z", "a"},
			wantError: true,
		},
	}
	for _, c := range cases {
		x, y, z = "", "", ""
		err := p.parse(c.args)
		if (err != nil) != c.wantError {
			errorString := "nil"
			if c.wantError {
				errorString = "non-nil"
			}
			t.Errorf("parse(%v) == %v, want %v", c.args, err, errorString)
		}
		if c.wantError {
			continue
		}
		if !(x == c.wantX && y == c.wantY && z == c.wantZ) {
			t.Errorf("parse(%v) set x,y,z to %v,%v,%v, want %v,%v,%v",
				c.args, x, y, z, c.wantX, c.wantY, c.wantZ)
		}
	}
}
