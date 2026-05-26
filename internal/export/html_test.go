package export

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/confsnap/internal/baseline"
)

func makeDriftResultForHTML(path string, status baseline.DriftStatus, baseHash, curHash string) baseline.DriftResult {
	return baseline.DriftResult{
		Path:         path,
		Status:       status,
		BaselineHash: baseHash,
		CurrentHash:  curHash,
	}
}

func TestDriftToHTML_WritesHeader(t *testing.T) {
	var buf bytes.Buffer
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	err := DriftToHTML(&buf, nil, "v1.0", ts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "<title>Drift Report: v1.0</title>") {
		t.Error("expected title with label")
	}
	if !strings.Contains(out, "<th>Path</th>") {
		t.Error("expected table header with Path")
	}
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Error("expected captured timestamp in output")
	}
}

func TestDriftToHTML_SortsRows(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForHTML("/etc/z.conf", baseline.StatusUnchanged, "aaa", "aaa"),
		makeDriftResultForHTML("/etc/a.conf", baseline.StatusModified, "bbb", "ccc"),
	}
	var buf bytes.Buffer
	err := DriftToHTML(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	idxA := strings.Index(out, "/etc/a.conf")
	idxZ := strings.Index(out, "/etc/z.conf")
	if idxA == -1 || idxZ == -1 {
		t.Fatal("expected both paths in output")
	}
	if idxA > idxZ {
		t.Error("expected /etc/a.conf before /etc/z.conf")
	}
}

func TestDriftToHTML_EmptyHashRenderedAsDash(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForHTML("/etc/new.conf", baseline.StatusAdded, "", "deadbeef"),
	}
	var buf bytes.Buffer
	err := DriftToHTML(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "><td>-</td>") {
		t.Error("expected empty baseline hash rendered as dash")
	}
}

func TestDriftToHTML_StatusCSSClass(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForHTML("/etc/removed.conf", baseline.StatusRemoved, "abc", ""),
	}
	var buf bytes.Buffer
	err := DriftToHTML(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `class="removed"`) {
		t.Error("expected removed CSS class on status cell")
	}
}

// TestDriftToHTML_AllStatusCSSClasses verifies that each drift status value
// produces the expected CSS class on the status table cell.
func TestDriftToHTML_AllStatusCSSClasses(t *testing.T) {
	cases := []struct {
		status    baseline.DriftStatus
		wantClass string
	}{
		{baseline.StatusUnchanged, "unchanged"},
		{baseline.StatusModified, "modified"},
		{baseline.StatusAdded, "added"},
		{baseline.StatusRemoved, "removed"},
	}
	for _, tc := range cases {
		t.Run(string(tc.status), func(t *testing.T) {
			results := []baseline.DriftResult{
				makeDriftResultForHTML("/etc/test.conf", tc.status, "abc", "abc"),
			}
			var buf bytes.Buffer
			if err := DriftToHTML(&buf, results, "test", time.Now()); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(buf.String(), `class="`+tc.wantClass+`"`) {
				t.Errorf("expected CSS class %q for status %q", tc.wantClass, tc.status)
			}
		})
	}
}
