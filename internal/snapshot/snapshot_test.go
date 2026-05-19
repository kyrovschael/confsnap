package snapshot_test

import (
	"os"
	"testing"

	"github.com/confsnap/confsnap/internal/snapshot"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "confsnap-*.conf")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestCapture_BasicFields(t *testing.T) {
	path := writeTempFile(t, "key=value\n")
	snap, err := snapshot.Capture(path, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Path != path {
		t.Errorf("expected path %q, got %q", path, snap.Path)
	}
	if snap.Hash == "" {
		t.Error("expected non-empty hash")
	}
	if snap.Size != int64(len("key=value\n")) {
		t.Errorf("expected size %d, got %d", len("key=value\n"), snap.Size)
	}
	if snap.Content != nil {
		t.Error("expected nil content when includeContent=false")
	}
}

func TestCapture_WithContent(t *testing.T) {
	const body = "port=8080\nhost=localhost\n"
	path := writeTempFile(t, body)
	snap, err := snapshot.Capture(path, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(snap.Content) != body {
		t.Errorf("content mismatch: got %q", snap.Content)
	}
}

func TestCapture_HashConsistency(t *testing.T) {
	path := writeTempFile(t, "stable=true\n")
	s1, _ := snapshot.Capture(path, false)
	s2, _ := snapshot.Capture(path, false)
	if s1.Hash != s2.Hash {
		t.Errorf("hashes differ for same file: %q vs %q", s1.Hash, s2.Hash)
	}
}

func TestCapture_MissingFile(t *testing.T) {
	_, err := snapshot.Capture("/nonexistent/path/config.conf", false)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
