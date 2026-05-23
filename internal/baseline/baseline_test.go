package baseline

import (
	"os"
	"testing"
	"time"
)

func setupBaselineDir(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp) // not used directly, but isolates env
	// Override baselineDir for tests via package-level var would be ideal;
	// here we rely on t.Chdir to isolate the working directory.
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
}

func makeBaseline(name string) *Baseline {
	return &Baseline{
		Name:      name,
		CreatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Label:     "v1.0",
		Files:     []string{"/etc/nginx/nginx.conf", "/etc/hosts"},
		Hashes: map[string]string{
			"/etc/nginx/nginx.conf": "abc123",
			"/etc/hosts":           "def456",
		},
	}
}

func TestSave_And_Load(t *testing.T) {
	setupBaselineDir(t)

	b := makeBaseline("production")
	if err := Save(b); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load("production")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.Name != b.Name {
		t.Errorf("Name: got %q, want %q", got.Name, b.Name)
	}
	if got.Label != b.Label {
		t.Errorf("Label: got %q, want %q", got.Label, b.Label)
	}
	if len(got.Files) != len(b.Files) {
		t.Errorf("Files len: got %d, want %d", len(got.Files), len(b.Files))
	}
	if got.Hashes["/etc/hosts"] != "def456" {
		t.Errorf("Hash mismatch for /etc/hosts")
	}
}

func TestLoad_NotFound(t *testing.T) {
	setupBaselineDir(t)

	_, err := Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing baseline, got nil")
	}
}

func TestList_Empty(t *testing.T) {
	setupBaselineDir(t)

	names, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestList_MultipleBaselines(t *testing.T) {
	setupBaselineDir(t)

	for _, name := range []string{"prod", "staging", "dev"} {
		if err := Save(makeBaseline(name)); err != nil {
			t.Fatalf("Save %q: %v", name, err)
		}
	}

	names, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 baselines, got %d: %v", len(names), names)
	}
}
