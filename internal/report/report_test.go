package report_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/user/confsnap/internal/report"
	"github.com/user/confsnap/internal/snapshot"
)

func setupSnapshots(t *testing.T, labelA, labelB string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CONFSNAP_DIR", dir)

	fileA := writeTempFile(t, "value=old\n")
	fileB := writeTempFile(t, "value=new\n")

	snapsA := []*snapshot.Snapshot{{Path: fileA, Hash: "aaa", Size: 9}}
	snapsB := []*snapshot.Snapshot{{Path: fileB, Hash: "bbb", Size: 9}}

	if err := snapshot.Save(labelA, snapsA); err != nil {
		t.Fatalf("save A: %v", err)
	}
	if err := snapshot.Save(labelB, snapsB); err != nil {
		t.Fatalf("save B: %v", err)
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "confsnap-*.conf")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestReport_Summary(t *testing.T) {
	setupSnapshots(t, "v1", "v2")
	r, err := report.New("v1", "v2")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if r.Summary() == "" {
		t.Error("expected non-empty summary")
	}
}

func TestReport_WriteText(t *testing.T) {
	setupSnapshots(t, "v1", "v2")
	r, err := report.New("v1", "v2")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatText); err != nil {
		t.Fatalf("Write text: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty text output")
	}
}

func TestReport_WriteJSON(t *testing.T) {
	setupSnapshots(t, "v1", "v2")
	r, err := report.New("v1", "v2")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatJSON); err != nil {
		t.Fatalf("Write JSON: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["label_a"] != "v1" || out["label_b"] != "v2" {
		t.Errorf("unexpected labels in JSON: %v", out)
	}
}
