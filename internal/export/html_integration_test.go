package export_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/confsnap/internal/baseline"
	"github.com/yourusername/confsnap/internal/export"
)

// TestDriftToHTML_FullRoundTrip verifies that all four drift statuses are
// represented correctly in a complete HTML document.
func TestDriftToHTML_FullRoundTrip(t *testing.T) {
	results := []baseline.DriftResult{
		{Path: "/etc/hosts", Status: baseline.StatusUnchanged, BaselineHash: "aaa", CurrentHash: "aaa"},
		{Path: "/etc/nginx.conf", Status: baseline.StatusModified, BaselineHash: "bbb", CurrentHash: "ccc"},
		{Path: "/etc/removed.conf", Status: baseline.StatusRemoved, BaselineHash: "ddd", CurrentHash: ""},
		{Path: "/etc/added.conf", Status: baseline.StatusAdded, BaselineHash: "", CurrentHash: "eee"},
	}

	var buf bytes.Buffer
	ts := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	err := export.DriftToHTML(&buf, results, "prod-v2", ts)
	if err != nil {
		t.Fatalf("DriftToHTML returned error: %v", err)
	}

	out := buf.String()

	expected := []string{
		"<!DOCTYPE html>",
		"prod-v2",
		"2024-01-15T09:00:00Z",
		"/etc/hosts",
		"/etc/nginx.conf",
		"/etc/removed.conf",
		"/etc/added.conf",
		`class="unchanged"`,
		`class="modified"`,
		`class="removed"`,
		`class="added"`,
		"</html>",
	}

	for _, want := range expected {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}
