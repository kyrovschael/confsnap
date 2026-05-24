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

// SlackPayload represents a Slack incoming webhook message.
type SlackPayload struct {
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

// SlackAttachment represents a single Slack message attachment.
type SlackAttachment struct {
	Color  string       `json:"color"`
	Fields []SlackField `json:"fields"`
}

// SlackField is a key/value pair inside an attachment.
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// DriftToSlack writes a Slack-compatible JSON webhook payload summarising
// the drift results to w.
func DriftToSlack(w io.Writer, results []baseline.DriftResult, label string, ts time.Time) error {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	var changed, removed, unchanged int
	for _, r := range results {
		switch r.Status {
		case baseline.StatusModified:
			changed++
		case baseline.StatusRemoved:
			removed++
		default:
			unchanged++
		}
	}

	header := fmt.Sprintf("*confsnap drift report* — baseline: `%s` — %s",
		label, ts.UTC().Format(time.RFC3339))

	color := "good"
	if changed > 0 || removed > 0 {
		color = "danger"
	}

	var fields []SlackField
	for _, r := range results {
		hash := r.CurrentHash
		if hash == "" {
			hash = "-"
		}
		fields = append(fields, SlackField{
			Title: r.Path,
			Value: fmt.Sprintf("status: %s | hash: %s", r.Status, hash),
			Short: false,
		})
	}

	payload := SlackPayload{
		Text: header,
		Attachments: []SlackAttachment{
			{
				Color:  color,
				Fields: fields,
			},
		},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = bytes.NewBuffer(nil) // keep import
	return enc.Encode(payload)
}
