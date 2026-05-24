package export

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/confsnap/internal/baseline"
)

func makeDriftResultForMD(path, status, expected, actual string) baseline.DriftResult {
	return baseline.DriftResult{
		Path:         path,
		Status:       baseline.DriftStatus(status),
		ExpectedHash: expected,
		ActualHash:   actual,
	}
}

func TestDriftToMarkdown_WritesHeader(t *testing.T) {
	var buf bytes.Buffer
	at := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	err := DriftToMarkdown(&buf, nil, "prod-v1", at)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "# Drift Report") {
		t.Error("missing markdown heading")
	}
	if !strings.Contains(out, "prod-v1") {
		t.Error("missing baseline name")
	}
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Error("missing timestamp")
	}
	if !strings.Contains(out, "| File | Status |") {
		t.Error("missing table header")
	}
}

func TestDriftToMarkdown_SortsRows(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForMD("/etc/zz.conf", "unchanged", "aaa", "aaa"),
		makeDriftResultForMD("/etc/aa.conf", "modified", "bbb", "ccc"),
	}
	var buf bytes.Buffer
	err := DriftToMarkdown(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	aaIdx := strings.Index(out, "aa.conf")
	zzIdx := strings.Index(out, "zz.conf")
	if aaIdx > zzIdx {
		t.Error("rows not sorted alphabetically by path")
	}
}

func TestDriftToMarkdown_EmptyHashRenderedAsDash(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForMD("/etc/missing.conf", "removed", "abc123", ""),
	}
	var buf bytes.Buffer
	err := DriftToMarkdown(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "—") {
		t.Error("expected em-dash for missing hash")
	}
}
