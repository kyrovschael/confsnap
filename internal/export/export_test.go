package export_test

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/user/confsnap/internal/diff"
	"github.com/user/confsnap/internal/export"
)

func makeResult(file, status, oldHash, newHash string, oldSize, newSize int64) diff.Result {
	var s diff.Status
	switch status {
	case "added":
		s = diff.Added
	case "removed":
		s = diff.Removed
	case "changed":
		s = diff.Changed
	default:
		s = diff.Unchanged
	}
	return diff.Result{
		File:    file,
		Status:  s,
		OldHash: oldHash,
		NewHash: newHash,
		OldSize: oldSize,
		NewSize: newSize,
	}
}

func TestToCSV_WritesHeader(t *testing.T) {
	var buf bytes.Buffer
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	if err := export.ToCSV(&buf, nil, ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid csv: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 row (header), got %d", len(records))
	}
	if records[0][0] != "file" {
		t.Errorf("expected first header column 'file', got %q", records[0][0])
	}
}

func TestToCSV_SortsRows(t *testing.T) {
	results := []diff.Result{
		makeResult("/etc/zsh/zshrc", "unchanged", "aaa", "aaa", 100, 100),
		makeResult("/etc/apt/sources.list", "changed", "bbb", "ccc", 200, 210),
		makeResult("/etc/hosts", "added", "", "ddd", 0, 50),
	}
	ts := time.Now()

	var buf bytes.Buffer
	if err := export.ToCSV(&buf, results, ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(&buf)
	records, _ := r.ReadAll()
	// records[0] is header
	if records[1][0] != "/etc/apt/sources.list" {
		t.Errorf("expected first data row /etc/apt/sources.list, got %q", records[1][0])
	}
	if records[3][0] != "/etc/zsh/zshrc" {
		t.Errorf("expected last data row /etc/zsh/zshrc, got %q", records[3][0])
	}
}

func TestToCSV_TimestampFormat(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC)
	results := []diff.Result{makeResult("/etc/hosts", "unchanged", "abc", "abc", 10, 10)}

	var buf bytes.Buffer
	if err := export.ToCSV(&buf, results, ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "2024-06-01T12:30:00Z") {
		t.Errorf("expected RFC3339 timestamp in output, got:\n%s", buf.String())
	}
}
