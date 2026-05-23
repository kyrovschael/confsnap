package schedule

import (
	"os"
	"testing"
	"time"
)

func overrideDir(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
}

func TestSaveAndLoad(t *testing.T) {
	overrideDir(t)
	s := &Schedule{
		Entries: []Entry{
			{Label: "nginx", Files: []string{"/etc/nginx/nginx.conf"}, Interval: time.Hour},
		},
	}
	if err := Save(s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Label != "nginx" {
		t.Errorf("expected label nginx, got %s", loaded.Entries[0].Label)
	}
}

func TestLoad_NotFound(t *testing.T) {
	overrideDir(t)
	s, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Entries) != 0 {
		t.Errorf("expected empty schedule")
	}
}

func TestAdd_UpdatesExisting(t *testing.T) {
	s := &Schedule{}
	s.Add(Entry{Label: "app", Interval: time.Minute})
	s.Add(Entry{Label: "app", Interval: time.Hour})
	if len(s.Entries) != 1 {
		t.Errorf("expected 1 entry after update, got %d", len(s.Entries))
	}
	if s.Entries[0].Interval != time.Hour {
		t.Errorf("expected updated interval")
	}
}

func TestRemove(t *testing.T) {
	s := &Schedule{}
	s.Add(Entry{Label: "app"})
	if !s.Remove("app") {
		t.Error("expected Remove to return true")
	}
	if len(s.Entries) != 0 {
		t.Error("expected empty entries after remove")
	}
	if s.Remove("missing") {
		t.Error("expected Remove to return false for missing label")
	}
}

func TestDue(t *testing.T) {
	now := time.Now()
	s := &Schedule{
		Entries: []Entry{
			{Label: "fresh", Interval: time.Hour, LastRun: now.Add(-30 * time.Minute)},
			{Label: "overdue", Interval: time.Hour, LastRun: now.Add(-2 * time.Hour)},
			{Label: "never", Interval: time.Hour},
		},
	}
	due := s.Due(now)
	if len(due) != 2 {
		t.Fatalf("expected 2 due entries, got %d", len(due))
	}
	_ = os.Getenv("HOME") // suppress unused import warning
}
