package cmd

import (
	"strings"
	"testing"
	"time"
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

func TestFlagsUsage(t *testing.T) {
	c := New("command", nil)
	c.Flag("-x", new(bool), "")
	c.Flag("-y", new(bool), "")
	got := c.flagsUsage()
	want := "[OPTION]..."
	if got != want {
		t.Errorf("usage returned %v, want %v", got, want)
	}
}

func TestParse(t *testing.T) {
	var (
		size     int
		timeout  time.Duration
		v        bool
		percent  float64
		count    int
		distance int
		name     string
	)
	c := New("command", nil)
	c.Bytes("-m --max-size", &size, "SIZE", "")
	c.Duration("-timeout", &timeout, "D", "")
	c.Flag("-v", &v, "")
	c.Float("--percent", &percent, "P", "")
	c.Int("--count", &count, "N", "")
	c.Metric("-d", &distance, "DISTANCE", "")
	c.String("-name", &name, "NAME", "")

	args := "-m 2k -timeout 5m -v --percent 99.5 --count 7 -d 150G -name moon"
	c.parse(strings.Split(args, " "))
	wantSize := 2048
	wantTimeout := 5 * time.Minute
	wantV := true
	wantPercent := 99.5
	wantCount := 7
	wantDistance := 150 * 1000000000
	wantName := "moon"
	if size != wantSize {
		t.Errorf("parse set size = %v, want %v", size, wantSize)
	}
	if timeout != wantTimeout {
		t.Errorf("parse set timeout = %v, want %v", timeout, wantTimeout)
	}
	if v != wantV {
		t.Errorf("parse set v = %v, want %v", v, wantV)
	}
	if percent != wantPercent {
		t.Errorf("parse set percent = %v, want %v", percent, wantPercent)
	}
	if count != wantCount {
		t.Errorf("parse set count = %v, want %v", count, wantCount)
	}
	if distance != wantDistance {
		t.Errorf("parse set distance = %v, want %v", distance, wantDistance)
	}
	if name != wantName {
		t.Errorf("parse set name = %v, want %v", name, wantName)
	}
}
