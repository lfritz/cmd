package cmd

import (
	"os"
	"reflect"
	"testing"
)

func TestWrapParagraphs(t *testing.T) {
	text := `
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.

Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
`
	columns := 30
	got := wrapParagraphs(text, columns)
	want := `Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor
incididunt ut labore et dolore
magna aliqua.

Ut enim ad minim veniam, quis
nostrud exercitation ullamco
laboris nisi ut aliquip ex ea
commodo consequat. Duis aute
irure dolor in reprehenderit
in voluptate velit esse cillum
dolore eu fugiat nulla
pariatur.

Excepteur sint occaecat
cupidatat non proident, sunt
in culpa qui officia deserunt
mollit anim id est laborum.
`
	if got != want {
		t.Errorf("wrapParagraphs(%v) == `%v`, want `%v`", columns, got, want)
	}
}

func TestWrapText(t *testing.T) {
	text := `
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
`
	columns := 30
	got := wrapText(text, columns)
	want := []string{
		"Lorem ipsum dolor sit amet,",
		"consectetur adipiscing elit,",
		"sed do eiusmod tempor",
		"incididunt ut labore et dolore",
		"magna aliqua. Ut enim ad minim",
		"veniam, quis nostrud",
		"exercitation ullamco laboris",
		"nisi ut aliquip ex ea commodo",
		"consequat. Duis aute irure",
		"dolor in reprehenderit in",
		"voluptate velit esse cillum",
		"dolore eu fugiat nulla",
		"pariatur. Excepteur sint",
		"occaecat cupidatat non",
		"proident, sunt in culpa qui",
		"officia deserunt mollit anim",
		"id est laborum.",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wrapText(text, %v) == `%v`, want `%v`", columns, got, want)
	}
}

func TestFormatTerms(t *testing.T) {
	def := &definition{
		terms: []string{"-H", "--dereference-command-line"},
	}
	cases := []struct {
		columns int
		want    termLines
	}{
		{
			columns: 30,
			want: termLines{
				separate: nil,
				inline:   "-H, --dereference-command-line",
			},
		},
		{
			columns: 26,
			want: termLines{
				separate: []string{"-H"},
				inline:   "--dereference-command-line",
			},
		},
		{
			columns: 25,
			want: termLines{
				separate: []string{"-H", "--dereference-command-line"},
				inline:   "",
			},
		},
	}
	for _, c := range cases {
		got := def.formatTerms(c.columns)
		if !reflect.DeepEqual(got.separate, c.want.separate) {
			t.Errorf("formatTerms returned separate == %v, want %v",
				got.separate, c.want.separate)
		}
		if got.inline != c.want.inline {
			t.Errorf("formatTerms returned inline == %v, want %v",
				got.inline, c.want.inline)
		}
	}
}

var optionsList = &definitionList{
	title: "Options",
	definitions: []*definition{
		{
			terms: []string{"-a", "--all"},
			text:  "do not ignore entries starting with .",
		},
		{
			terms: []string{"-H", "--dereference-command-line"},
			text:  "follow symbolic links listed on the command line",
		},
	},
}
var commandsList = &definitionList{
	title: "Commands",
	definitions: []*definition{
		{
			terms: []string{"go"},
			text:  "Go go go!",
		},
	},
}

// func (d *definitionList) format(columns int) string {
func TestDefinitionsList(t *testing.T) {
	columns := 30
	got := optionsList.format(columns)
	want := `Options:
  -a, --all  do not ignore
             entries starting
             with .
  -H
  --dereference-command-line
             follow symbolic
             links listed on
             the command line
`
	if got != want {
		t.Errorf("optionsList.format(%v) == `%v`, want `%v`", columns, got, want)
	}
}

func TestFormatHelp(t *testing.T) {
	usage := "Usage: ls [OPTION]... [FILE]..."
	summary := `
List information about the FILEs (the current directory by default).
Sort entries alphabetically if none of -cftuvSUX nor --sort is specified.`
	details := `
The SIZE argument is an integer and optional unit (example: 10K is 10*1024).
Units are K,M,G,T,P,E,Z,Y (powers of 1024) or KB,MB,... (powers of 1000).

The TIME_STYLE argument can be full-iso, long-iso, iso, locale, or +FORMAT.
FORMAT is interpreted like in date(1).  If FORMAT is FORMAT1<newline>FORMAT2,
then FORMAT1 applies to non-recent files and FORMAT2 to recent files.
TIME_STYLE prefixed with 'posix-' takes effect only outside the POSIX locale.
Also the TIME_STYLE environment variable sets the default style to use.
`
	defs := []*definitionList{
		optionsList,
		commandsList,
	}
	want := `Usage: ls [OPTION]... [FILE]...

List information about the FILEs (the current directory by default). Sort entries
alphabetically if none of -cftuvSUX nor --sort is specified.

Options:
  -a, --all  do not ignore entries starting with .
  -H
  --dereference-command-line
             follow symbolic links listed on the command line

Commands:
  go  Go go go!

The SIZE argument is an integer and optional unit (example: 10K is 10*1024). Units are
K,M,G,T,P,E,Z,Y (powers of 1024) or KB,MB,... (powers of 1000).

The TIME_STYLE argument can be full-iso, long-iso, iso, locale, or +FORMAT. FORMAT is
interpreted like in date(1). If FORMAT is FORMAT1<newline>FORMAT2, then FORMAT1 applies to
non-recent files and FORMAT2 to recent files. TIME_STYLE prefixed with 'posix-' takes
effect only outside the POSIX locale. Also the TIME_STYLE environment variable sets the
default style to use.
`
	os.Setenv("COLUMNS", "90")
	got := formatHelp(usage, summary, details, defs)
	if got != want {
		t.Errorf("formatHelp returned `%v`, want `%v`", got, want)
	}
}
