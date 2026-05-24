package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/confsnap/internal/baseline"
)

func makeDriftResultForSlack(path, status, hash string) baseline.DriftResult {
	return baseline.DriftResult{
		Path:        path,
		Status:      baseline.DriftStatus(status),
		CurrentHash: hash,
	}
}

func TestDriftToSlack_WritesHeader(t *testing.T) {
	var buf bytes.Buffer
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	err := DriftToSlack(&buf, nil, "prod-v1", ts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "prod-v1") {
		t.Error("expected baseline label in output")
	}
	if !strings.Contains(buf.String(), "2024-06-01T12:00:00Z") {
		t.Error("expected timestamp in output")
	}
}

func TestDriftToSlack_SortsRows(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForSlack("/etc/z.conf", "unchanged", "aaa"),
		makeDriftResultForSlack("/etc/a.conf", "modified", "bbb"),
	}
	var buf bytes.Buffer
	err := DriftToSlack(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	idxA := strings.Index(buf.String(), "/etc/a.conf")
	idxZ := strings.Index(buf.String(), "/etc/z.conf")
	if idxA > idxZ {
		t.Error("expected rows sorted by path")
	}
}

func TestDriftToSlack_EmptyHashRenderedAsDash(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForSlack("/etc/missing.conf", "removed", ""),
	}
	var buf bytes.Buffer
	err := DriftToSlack(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hash: -") {
		t.Error("expected empty hash rendered as dash")
	}
}

func TestDriftToSlack_ColorDangerOnDrift(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForSlack("/etc/nginx.conf", "modified", "abc123"),
	}
	var buf bytes.Buffer
	err := DriftToSlack(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload SlackPayload
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(payload.Attachments) == 0 {
		t.Fatal("expected at least one attachment")
	}
	if payload.Attachments[0].Color != "danger" {
		t.Errorf("expected color 'danger', got %q", payload.Attachments[0].Color)
	}
}

func TestDriftToSlack_ColorGoodWhenClean(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForSlack("/etc/nginx.conf", "unchanged", "abc123"),
	}
	var buf bytes.Buffer
	err := DriftToSlack(&buf, results, "test", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload SlackPayload
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Attachments[0].Color != "good" {
		t.Errorf("expected color 'good', got %q", payload.Attachments[0].Color)
	}
}
