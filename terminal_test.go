package cmd

import (
	"os"
	"strconv"
	"testing"
)

func TestTerminalColumns(t *testing.T) {
	want := 99
	os.Setenv("COLUMNS", strconv.Itoa(want))
	got := terminalColumns()
	if got != want {
		t.Errorf("terminalColumns() == %v, want %v (using env var)", got, want)
	}

	os.Unsetenv("COLUMNS")
	_ = terminalColumns() // just make sure it doesn't crash
}
