package export

import (
	"strings"
	"testing"
	"time"

	"github.com/nicholasgasior/confsnap/internal/baseline"
)

func makeDriftResultForPDF(path, hash string, status baseline.DriftStatus) baseline.DriftResult {
	return baseline.DriftResult{
		Path:        path,
		CurrentHash: hash,
		Status:      status,
	}
}

func TestDriftToPDF_WritesHeader(t *testing.T) {
	var buf strings.Builder
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	err := DriftToPDF(&buf, nil, ts, "prod-v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CONFSNAP DRIFT REPORT") {
		t.Error("missing report title")
	}
	if !strings.Contains(out, "prod-v1") {
		t.Error("missing baseline label")
	}
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Error("missing captured timestamp")
	}
}

func TestDriftToPDF_SortsRows(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForPDF("/etc/z.conf", "aaa", baseline.StatusUnchanged),
		makeDriftResultForPDF("/etc/a.conf", "bbb", baseline.StatusModified),
	}
	var buf strings.Builder
	ts := time.Now()
	if err := DriftToPDF(&buf, results, ts, "test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	idxA := strings.Index(out, "/etc/a.conf")
	idxZ := strings.Index(out, "/etc/z.conf")
	if idxA > idxZ {
		t.Error("rows not sorted alphabetically by path")
	}
}

func TestDriftToPDF_EmptyHashRenderedAsDash(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForPDF("/etc/missing.conf", "", baseline.StatusRemoved),
	}
	var buf strings.Builder
	if err := DriftToPDF(&buf, results, time.Now(), "base"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "-") {
		t.Error("expected dash for empty hash")
	}
}

func TestDriftToPDF_ContainsColumnHeaders(t *testing.T) {
	var buf strings.Builder
	if err := DriftToPDF(&buf, nil, time.Now(), "base"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"PATH", "STATUS", "HASH"} {
		if !strings.Contains(out, col) {
			t.Errorf("missing column header: %s", col)
		}
	}
}
