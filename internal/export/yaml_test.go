package export

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/nicholasgasior/confsnap/internal/baseline"
)

func makeDriftResultForYAML(path, baseHash, curHash string, status baseline.DriftStatus) baseline.DriftResult {
	return baseline.DriftResult{
		Path:         path,
		BaselineHash: baseHash,
		CurrentHash:  curHash,
		Status:       status,
	}
}

func TestDriftToYAML_WritesHeader(t *testing.T) {
	var buf bytes.Buffer
	t0 := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	err := DriftToYAML(&buf, nil, "prod-v1", t0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "baseline: \"prod-v1\"") {
		t.Errorf("expected baseline label in output, got:\n%s", out)
	}
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Errorf("expected captured_at timestamp in output, got:\n%s", out)
	}
}

func TestDriftToYAML_SortsRows(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForYAML("/etc/zz.conf", "aaa", "bbb", baseline.StatusModified),
		makeDriftResultForYAML("/etc/aa.conf", "ccc", "ccc", baseline.StatusUnchanged),
	}
	var buf bytes.Buffer
	_ = DriftToYAML(&buf, results, "test", time.Now())
	out := buf.String()
	idxAA := strings.Index(out, "aa.conf")
	idxZZ := strings.Index(out, "zz.conf")
	if idxAA > idxZZ {
		t.Errorf("expected aa.conf before zz.conf in YAML output")
	}
}

func TestDriftToYAML_EmptyHashRenderedAsDash(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForYAML("/etc/missing.conf", "abc123", "", baseline.StatusRemoved),
	}
	var buf bytes.Buffer
	_ = DriftToYAML(&buf, results, "test", time.Now())
	out := buf.String()
	if !strings.Contains(out, "current_hash: \"-\"") {
		t.Errorf("expected dash for empty current hash, got:\n%s", out)
	}
}

func TestDriftToYAML_SummaryCounts(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForYAML("/a", "h1", "h1", baseline.StatusUnchanged),
		makeDriftResultForYAML("/b", "h2", "h3", baseline.StatusModified),
		makeDriftResultForYAML("/c", "h4", "", baseline.StatusRemoved),
	}
	var buf bytes.Buffer
	_ = DriftToYAML(&buf, results, "test", time.Now())
	out := buf.String()
	if !strings.Contains(out, "unchanged: 1") {
		t.Errorf("expected unchanged: 1, got:\n%s", out)
	}
	if !strings.Contains(out, "modified: 1") {
		t.Errorf("expected modified: 1, got:\n%s", out)
	}
	if !strings.Contains(out, "removed: 1") {
		t.Errorf("expected removed: 1, got:\n%s", out)
	}
}
