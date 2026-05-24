package export

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/confsnap/internal/baseline"
)

func makeDriftResultForEmail(path string, status baseline.DriftStatus, hash string) baseline.DriftResult {
	return baseline.DriftResult{
		Path:        path,
		Status:      status,
		CurrentHash: hash,
	}
}

func TestDriftToEmail_WritesSubjectAndHeaders(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForEmail("/etc/hosts", baseline.StatusUnchanged, "abc123"),
	}
	var buf bytes.Buffer
	if err := DriftToEmail(&buf, results, "deploy-42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Subject: confsnap drift report") {
		t.Error("expected Subject header in output")
	}
	if !strings.Contains(out, "deploy-42") {
		t.Error("expected snapshot label in output")
	}
	if !strings.Contains(out, "MIME-Version: 1.0") {
		t.Error("expected MIME-Version header")
	}
}

func TestDriftToEmail_SortsRows(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForEmail("/etc/zshrc", baseline.StatusModified, "zzz"),
		makeDriftResultForEmail("/etc/apt/sources.list", baseline.StatusAdded, "aaa"),
		makeDriftResultForEmail("/etc/hosts", baseline.StatusUnchanged, "bbb"),
	}
	var buf bytes.Buffer
	if err := DriftToEmail(&buf, results, "test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	aptIdx := strings.Index(out, "/etc/apt/sources.list")
	hostsIdx := strings.Index(out, "/etc/hosts")
	zshIdx := strings.Index(out, "/etc/zshrc")
	if aptIdx > hostsIdx || hostsIdx > zshIdx {
		t.Error("rows are not sorted by path")
	}
}

func TestDriftToEmail_EmptyHashRenderedAsDash(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForEmail("/etc/removed", baseline.StatusRemoved, ""),
	}
	var buf bytes.Buffer
	if err := DriftToEmail(&buf, results, "test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "-") {
		t.Error("expected dash placeholder for empty hash")
	}
}

func TestDriftToEmail_SummaryCounts(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForEmail("/a", baseline.StatusAdded, "h1"),
		makeDriftResultForEmail("/b", baseline.StatusRemoved, ""),
		makeDriftResultForEmail("/c", baseline.StatusModified, "h2"),
		makeDriftResultForEmail("/d", baseline.StatusUnchanged, "h3"),
	}
	var buf bytes.Buffer
	if err := DriftToEmail(&buf, results, "test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Added    : 1", "Removed  : 1", "Modified : 1", "Unchanged: 1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in summary", want)
		}
	}
}

func TestDriftToEmail_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	err := DriftToEmail(&buf, nil, "test")
	if err == nil {
		t.Error("expected error for empty results, got nil")
	}
}
