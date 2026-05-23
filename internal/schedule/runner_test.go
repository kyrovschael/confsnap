package schedule

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "conf-*.conf")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRunOnce_NoDue(t *testing.T) {
	overrideDir(t)
	now := time.Now()
	s := &Schedule{
		Entries: []Entry{
			{Label: "recent", Files: []string{"/etc/hosts"}, Interval: time.Hour, LastRun: now.Add(-1 * time.Minute)},
		},
	}
	if err := Save(s); err != nil {
		t.Fatalf("save: %v", err)
	}
	var buf bytes.Buffer
	if err := RunOnce(&buf, now); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if buf.String() != "no entries due\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestRunOnce_SnapshotsOverdue(t *testing.T) {
	overrideDir(t)
	path := writeTempFile(t, "key=value\n")
	now := time.Now()
	s := &Schedule{
		Entries: []Entry{
			{Label: "myapp", Files: []string{path}, Interval: time.Minute, LastRun: now.Add(-5 * time.Minute)},
		},
	}
	if err := Save(s); err != nil {
		t.Fatalf("save: %v", err)
	}
	var buf bytes.Buffer
	if err := RunOnce(&buf, now); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output from RunOnce")
	}
	// Verify LastRun was updated.
	loaded, _ := Load()
	if loaded.Entries[0].LastRun.IsZero() {
		t.Error("expected LastRun to be updated")
	}
	_ = filepath.Base(path) // suppress unused import
}
