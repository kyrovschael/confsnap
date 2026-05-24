package export

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/confsnap/internal/baseline"
)

func makeDriftResultForWebhook(file, hash, baselineHash string, status baseline.DriftStatus) baseline.DriftResult {
	return baseline.DriftResult{
		File:         file,
		Status:       status,
		CurrentHash:  hash,
		BaselineHash: baselineHash,
	}
}

func TestDriftToWebhook_WritesValidJSON(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForWebhook("/etc/hosts", "abc123", "abc123", baseline.StatusUnchanged),
		makeDriftResultForWebhook("/etc/nginx.conf", "def456", "aaa000", baseline.StatusModified),
	}

	var buf bytes.Buffer
	if err := DriftToWebhook(&buf, results, "v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload WebhookPayload
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestDriftToWebhook_SummaryCounts(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForWebhook("/a", "h1", "h1", baseline.StatusUnchanged),
		makeDriftResultForWebhook("/b", "h2", "h0", baseline.StatusModified),
		makeDriftResultForWebhook("/c", "", "h3", baseline.StatusRemoved),
	}

	data, err := WebhookBytes(results, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload WebhookPayload
	_ = json.Unmarshal(data, &payload)

	if payload.Summary.Total != 3 {
		t.Errorf("expected total 3, got %d", payload.Summary.Total)
	}
	if payload.Summary.Unchanged != 1 {
		t.Errorf("expected unchanged 1, got %d", payload.Summary.Unchanged)
	}
	if payload.Summary.Modified != 1 {
		t.Errorf("expected modified 1, got %d", payload.Summary.Modified)
	}
	if payload.Summary.Removed != 1 {
		t.Errorf("expected removed 1, got %d", payload.Summary.Removed)
	}
}

func TestDriftToWebhook_SortsRows(t *testing.T) {
	results := []baseline.DriftResult{
		makeDriftResultForWebhook("/z/last", "h1", "h1", baseline.StatusUnchanged),
		makeDriftResultForWebhook("/a/first", "h2", "h2", baseline.StatusUnchanged),
	}

	data, _ := WebhookBytes(results, "test")
	var payload WebhookPayload
	_ = json.Unmarshal(data, &payload)

	if payload.Results[0].File != "/a/first" {
		t.Errorf("expected first row /a/first, got %s", payload.Results[0].File)
	}
}

func TestDriftToWebhook_BaselineField(t *testing.T) {
	data, _ := WebhookBytes(nil, "release-42")
	var payload WebhookPayload
	_ = json.Unmarshal(data, &payload)

	if payload.Baseline != "release-42" {
		t.Errorf("expected baseline release-42, got %s", payload.Baseline)
	}
}
