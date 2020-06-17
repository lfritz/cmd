package cmd

import "testing"

func TestGroupUsage(t *testing.T) {
	g := NewGroup("service")
	g.Flag("-v", new(bool), "")
	start := g.Command("start", func() {})
	check := g.Group("check")
	database := check.Command("database", func() {})
	database.String("--host", new(string), "HOST", "")
	database.Int("--port", new(int), "PORT", "")

	want := "Usage: service [OPTION] GROUP | COMMAND"
	got := g.usage()
	if got != want {
		t.Errorf("g.usage() == `%s`, want `%s`", got, want)
	}

	want = "Usage: service start"
	got = start.usage()
	if got != want {
		t.Errorf("start.usage() == `%s`, want `%s`", got, want)
	}

	want = "Usage: service check COMMAND"
	got = check.usage()
	if got != want {
		t.Errorf("check.usage() == `%s`, want `%s`", got, want)
	}

	want = "Usage: service check database [OPTION]..."
	got = database.usage()
	if got != want {
		t.Errorf("database.usage() == `%s`, want `%s`", got, want)
	}
}

func TestGroupRun(t *testing.T) {
	var startCalled, databaseCalled bool
	runStart := func() {
		startCalled = true
	}
	runDatabase := func() {
		databaseCalled = true
	}
	g := NewGroup("service")
	g.Flag("-v", new(bool), "")
	g.Command("start", runStart)
	check := g.Group("check")
	database := check.Command("database", runDatabase)
	database.String("--host", new(string), "HOST", "")
	database.Int("--port", new(int), "PORT", "")

	g.run([]string{"start"}, false)
	if !startCalled {
		t.Errorf("Group.run didn't call expected function")
	}

	g.run([]string{"check", "database"}, false)
	if !databaseCalled {
		t.Errorf("Group.run didn't call expected function")
	}
}
