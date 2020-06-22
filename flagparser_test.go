package cmd

import (
	"reflect"
	"testing"
)

func TestSplitSpec(t *testing.T) {
	cases := []struct {
		spec      string
		want      []string
		wantError bool
	}{
		{"-v", []string{"-v"}, false},
		{"-v --verbose -d --debug", []string{"-v", "--verbose", "-d", "--debug"}, false},
		{"", nil, true},
		{"---verbose", nil, true},
		{"hello", nil, true},
	}
	for _, c := range cases {
		got, err := splitSpec(c.spec)
		if err != nil && !c.wantError {
			t.Errorf("splitSpec(%v) returned error", c.spec)
			continue
		}
		if err == nil && c.wantError {
			t.Errorf("splitSpec(%v) didn't return error", c.spec)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("splitSpec(%v) returned %v, want %v", c.spec, got, c.want)
		}
	}
}

func TestFlagsAndOptions(t *testing.T) {
	var (
		verbose, long bool
		color, shape  string
	)
	p := newFlagParser(true)
	p.addFlag("-v --verbose", &verbose)
	p.addFlag("--long", &long)
	p.addOption("--color", setter(&color))
	p.addOption("--shape", setter(&shape))
	args := []string{"-v", "--color=red", "--shape", "square", "hello", "world"}
	remaining, help, version, err := p.parse(args, true, true)
	wantRemaining := []string{"hello", "world"}
	wantVerbose := true
	wantLong := false
	wantColor := "red"
	wantShape := "square"
	if err != nil {
		t.Errorf("parse returned error: %v", err)
	}
	if !reflect.DeepEqual(remaining, wantRemaining) {
		t.Errorf("parse returned remaining = %v, want %v", remaining, wantRemaining)
	}
	if help {
		t.Error("parse returned help = true")
	}
	if version {
		t.Error("parse returned help = true")
	}
	if verbose != wantVerbose {
		t.Errorf("parse set verbose = %v, want %v", verbose, wantVerbose)
	}
	if long != wantLong {
		t.Errorf("parse set long = %v, want %v", long, wantLong)
	}
	if color != wantColor {
		t.Errorf("parse set color = %v, want %v", color, wantColor)
	}
	if shape != wantShape {
		t.Errorf("parse set shape = %v, want %v", shape, wantShape)
	}
}

func TestMixingFlagsAndArgs(t *testing.T) {
	var all, long bool
	p := newFlagParser(true)
	p.addFlag("-l", &long)
	p.addFlag("-a", &all)
	args := []string{"-a", "foo", "bar", "-l"}
	remaining, help, version, err := p.parse(args, false, true)
	wantRemaining := []string{"foo", "bar"}
	if err != nil {
		t.Errorf("parse returned error: %v", err)
	}
	if !reflect.DeepEqual(remaining, wantRemaining) {
		t.Errorf("parse returned remaining = %v, want %v", remaining, wantRemaining)
	}
	if help {
		t.Error("parse returned help = true")
	}
	if version {
		t.Error("parse returned help = true")
	}
	if !all {
		t.Errorf("parse didn't set all")
	}
	if !long {
		t.Errorf("parse didn't set long")
	}
}

func TestMergedFlags(t *testing.T) {
	var all, long bool
	p := newFlagParser(true)
	p.addFlag("-l", &long)
	p.addFlag("-a", &all)
	args := []string{"-la"}
	remaining, help, version, err := p.parse(args, false, true)
	wantRemaining := []string{}
	if err != nil {
		t.Errorf("parse returned error: %v", err)
	}
	if !reflect.DeepEqual(remaining, wantRemaining) {
		t.Errorf("parse returned remaining = %v, want %v", remaining, wantRemaining)
	}
	if help {
		t.Error("parse returned help = true")
	}
	if version {
		t.Error("parse returned help = true")
	}
	if !(all && long) {
		t.Error("parse didn't set flags")
	}
}

func setter(ptr *string) func(string, string) error {
	return func(_, value string) error {
		*ptr = value
		return nil
	}
}
