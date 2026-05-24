package export

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/confsnap/internal/baseline"
)

func makeDriftResult(file string, status baseline.DriftStatus) baseline.DriftResult {
	return baseline.DriftResult{File: file, Status: status}
}

func TestDriftToCSV_WritesHeader(t *testing.T) {
	var buf bytes.Buffer
	results := []baseline.DriftResult{}
	at := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	if err := DriftToCSV(&buf, results, at, "prod-v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line (header only), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "file,status,baseline,checked_at") {
		t.Errorf("unexpected header: %s", lines[0])
	}
}

func TestDriftToCSV_SortsRows(t *testing.T) {
	var buf bytes.Buffer
	at := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	results := []baseline.DriftResult{
		makeDriftResult("/etc/nginx.conf", baseline.StatusModified),
		makeDriftResult("/etc/hosts", baseline.StatusUnchanged),
		makeDriftResult("/etc/apt/sources.list", baseline.StatusRemoved),
	}

	if err := DriftToCSV(&buf, results, at, "base-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 3 data rows
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[1], "/etc/apt/sources.list") {
		t.Errorf("expected first data row to be /etc/apt/sources.list, got: %s", lines[1])
	}
	if !strings.HasPrefix(lines[2], "/etc/hosts") {
		t.Errorf("expected second data row to be /etc/hosts, got: %s", lines[2])
	}
}

func TestDriftToCSV_TimestampAndBaseline(t *testing.T) {
	var buf bytes.Buffer
	at := time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC)
	results := []baseline.DriftResult{
		makeDriftResult("/etc/hosts", baseline.StatusUnchanged),
	}

	if err := DriftToCSV(&buf, results, at, "release-42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := buf.String()
	if !strings.Contains(content, "2024-01-15T09:30:00Z") {
		t.Errorf("expected RFC3339 timestamp in output, got:\n%s", content)
	}
	if !strings.Contains(content, "release-42") {
		t.Errorf("expected baseline name in output, got:\n%s", content)
	}
}
