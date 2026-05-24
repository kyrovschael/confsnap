package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/confsnap/internal/baseline"
)

// WebhookPayload is the JSON structure posted to a webhook endpoint.
type WebhookPayload struct {
	GeneratedAt time.Time        `json:"generated_at"`
	Baseline    string           `json:"baseline"`
	Summary     WebhookSummary   `json:"summary"`
	Results     []WebhookEntry   `json:"results"`
}

// WebhookSummary contains aggregate counts for the payload.
type WebhookSummary struct {
	Total     int `json:"total"`
	Unchanged int `json:"unchanged"`
	Modified  int `json:"modified"`
	Removed   int `json:"removed"`
}

// WebhookEntry represents a single drift result in the payload.
type WebhookEntry struct {
	File     string `json:"file"`
	Status   string `json:"status"`
	Hash     string `json:"hash,omitempty"`
	Baseline string `json:"baseline_hash,omitempty"`
}

// DriftToWebhook encodes drift results as a JSON webhook payload written to w.
func DriftToWebhook(w io.Writer, results []baseline.DriftResult, baselineName string) error {
	sorted := make([]baseline.DriftResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].File < sorted[j].File
	})

	var summary WebhookSummary
	entries := make([]WebhookEntry, 0, len(sorted))

	for _, r := range sorted {
		summary.Total++
		switch r.Status {
		case baseline.StatusUnchanged:
			summary.Unchanged++
		case baseline.StatusModified:
			summary.Modified++
		case baseline.StatusRemoved:
			summary.Removed++
		}

		hash := r.CurrentHash
		if hash == "" {
			hash = ""
		}
		entries = append(entries, WebhookEntry{
			File:     r.File,
			Status:   r.Status.String(),
			Hash:     hash,
			Baseline: r.BaselineHash,
		})
	}

	payload := WebhookPayload{
		GeneratedAt: time.Now().UTC(),
		Baseline:    baselineName,
		Summary:     summary,
		Results:     entries,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		return fmt.Errorf("webhook encode: %w", err)
	}
	return nil
}

// WebhookBytes is a convenience wrapper that returns the payload as a byte slice.
func WebhookBytes(results []baseline.DriftResult, baselineName string) ([]byte, error) {
	var buf bytes.Buffer
	if err := DriftToWebhook(&buf, results, baselineName); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
