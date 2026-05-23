package watch

import (
	"os"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "watchtest-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestNew_InitialisesHashes(t *testing.T) {
	path := writeTempFile(t, "initial content")
	w, err := New([]string{path}, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if w.hashes[path] == "" {
		t.Error("expected non-empty initial hash")
	}
}

func TestPoll_NoChange(t *testing.T) {
	path := writeTempFile(t, "stable content")
	w, _ := New([]string{path}, time.Second)
	events, err := w.Poll()
	if err != nil {
		t.Fatalf("Poll: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestPoll_DetectsChange(t *testing.T) {
	path := writeTempFile(t, "original")
	w, _ := New([]string{path}, time.Second)

	if err := os.WriteFile(path, []byte("modified"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	events, err := w.Poll()
	if err != nil {
		t.Fatalf("Poll: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Path != path {
		t.Errorf("unexpected path: %s", events[0].Path)
	}
	if events[0].OldHash == events[0].NewHash {
		t.Error("expected old and new hashes to differ")
	}
}

func TestPoll_UpdatesHashAfterChange(t *testing.T) {
	path := writeTempFile(t, "v1")
	w, _ := New([]string{path}, time.Second)
	os.WriteFile(path, []byte("v2"), 0644)
	w.Poll() // consume change

	events, err := w.Poll()
	if err != nil {
		t.Fatalf("Poll: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events after hash update, got %d", len(events))
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := New([]string{"/nonexistent/path/file.conf"}, time.Second)
	if err == nil {
		t.Error("expected error for missing file")
	}
}
