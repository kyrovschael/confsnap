package baseline

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestCheckDrift_Unchanged(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "app.conf", "key=value\n")

	h, _ := hashFile(path)
	b := &Baseline{
		Name: "test", CreatedAt: time.Now(), Files: []string{path},
		Hashes: map[string]string{path: h},
	}

	drifts, err := CheckDrift(b)
	if err != nil {
		t.Fatalf("CheckDrift: %v", err)
	}
	if len(drifts) != 1 || drifts[0].Status != DriftUnchanged {
		t.Errorf("expected unchanged, got %v", drifts)
	}
}

func TestCheckDrift_Modified(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "app.conf", "key=value\n")

	b := &Baseline{
		Name: "test", CreatedAt: time.Now(), Files: []string{path},
		Hashes: map[string]string{path: "oldhash"},
	}

	drifts, err := CheckDrift(b)
	if err != nil {
		t.Fatalf("CheckDrift: %v", err)
	}
	if len(drifts) != 1 || drifts[0].Status != DriftModified {
		t.Errorf("expected modified, got %v", drifts)
	}
}

func TestCheckDrift_Removed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.conf")

	b := &Baseline{
		Name: "test", CreatedAt: time.Now(), Files: []string{path},
		Hashes: map[string]string{path: "somehash"},
	}

	drifts, err := CheckDrift(b)
	if err != nil {
		t.Fatalf("CheckDrift: %v", err)
	}
	if len(drifts) != 1 || drifts[0].Status != DriftRemoved {
		t.Errorf("expected removed, got %v", drifts)
	}
}

func TestDriftStatus_String(t *testing.T) {
	cases := []struct {
		s    DriftStatus
		want string
	}{
		{DriftUnchanged, "unchanged"},
		{DriftModified, "modified"},
		{DriftAdded, "added"},
		{DriftRemoved, "removed"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("DriftStatus(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}
