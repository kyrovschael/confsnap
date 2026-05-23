package report

import (
	"encoding/json"
	"io"
	"time"

	"github.com/user/confsnap/internal/diff"
)

type jsonEntry struct {
	Path   string `json:"path"`
	Status string `json:"status"`
	Diff   string `json:"diff,omitempty"`
}

type jsonReport struct {
	LabelA    string      `json:"label_a"`
	LabelB    string      `json:"label_b"`
	CreatedAt time.Time   `json:"created_at"`
	Summary   string      `json:"summary"`
	Entries   []jsonEntry `json:"entries"`
}

func writeJSON(w io.Writer, r *Report) error {
	entries := make([]jsonEntry, len(r.Entries))
	for i, e := range r.Entries {
		entries[i] = jsonEntry{
			Path:   e.Path,
			Status: statusString(e.Status),
			Diff:   e.Diff,
		}
	}
	payload := jsonReport{
		LabelA:    r.LabelA,
		LabelB:    r.LabelB,
		CreatedAt: r.CreatedAt,
		Summary:   r.Summary(),
		Entries:   entries,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func statusString(s diff.Status) string {
	switch s {
	case diff.StatusAdded:
		return "added"
	case diff.StatusRemoved:
		return "removed"
	case diff.StatusChanged:
		return "changed"
	default:
		return "unchanged"
	}
}
