package cmd

import (
	"reflect"
	"strings"
	"testing"
)

func TestFlagsParse(t *testing.T) {
	var (
		quiet bool
		name  string
	)
	flags := Flags{
		flags:   make(map[string]*bool),
		options: make(map[string]*option),
	}
	flags.Flag("--quiet", &quiet, "")
	flags.String("--name", &name, "NAME", "")

	cases := []struct {
		args      string
		err, help bool
		following []string
		name      string
	}{
		{
			args:      "--quiet --name Joe foo",
			name:      "Joe",
			following: []string{"foo"},
		},
		{
			args:      "--quiet --name=Joe foo",
			name:      "Joe",
			following: []string{"foo"},
		},
		{
			args: "--help",
			help: true,
		},
		{
			args: "--quiet --name",
			err:  true,
		},
	}

	for _, c := range cases {
		quiet = false
		name = ""

		args := strings.Split(c.args, " ")
		err, help, following := flags.parse(args)
		if c.err {
			if err == nil {
				t.Errorf("Flags.parse(%v) didn't return error", args)
			}
			continue
		}
		if err != nil {
			t.Errorf("Flags.parse returned error: %v", err)
		}
		if c.help {
			if !help {
				t.Errorf("Flags.parse(%v) didn't recognize help flag", args)
			}
			continue
		}
		if name != c.name {
			t.Errorf("Flags.parse(%v) set name to %v, want %v", args, name, c.name)
		}
		if !reflect.DeepEqual(following, c.following) {
			t.Errorf("Flags.parse(%v) returned %v, want %v", args, following, c.following)
		}
	}
}

func TestSplitSpec(t *testing.T) {
	spec := "-c --color"
	want := []string{"-c", "--color"}
	got, err := splitSpec(spec)
	if err != nil {
		t.Fatalf("splitSpec(%v) returned error: %v", spec, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("splitSpec(%v) == %v, want %v", spec, got, want)
	}

	for _, spec := range []string{
		"",
		"c",
		"-color shape",
		"---shape",
	} {
		_, err := splitSpec(spec)
		if err == nil {
			t.Errorf("splitSpec(%v) didn't return error", spec)
		}
	}
}
