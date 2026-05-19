package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SnapshotSet groups multiple file snapshots taken together under a label.
type SnapshotSet struct {
	Label     string      `json:"label"`
	CreatedAt time.Time   `json:"created_at"`
	Snapshots []*Snapshot `json:"snapshots"`
}

// Save persists a SnapshotSet as a JSON file inside dir.
// The filename is derived from the label and timestamp.
func Save(dir string, set *SnapshotSet) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating snapshot dir: %w", err)
	}

	ts := set.CreatedAt.Format("20060102T150405Z")
	filename := fmt.Sprintf("%s-%s.json", sanitizeLabel(set.Label), ts)
	outPath := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling snapshot set: %w", err)
	}

	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return "", fmt.Errorf("writing snapshot file: %w", err)
	}
	return outPath, nil
}

// Load reads a SnapshotSet from a JSON file.
func Load(path string) (*SnapshotSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot file %q: %w", path, err)
	}
	var set SnapshotSet
	if err := json.Unmarshal(data, &set); err != nil {
		return nil, fmt.Errorf("parsing snapshot file %q: %w", path, err)
	}
	return &set, nil
}

// sanitizeLabel replaces characters unsafe for filenames.
func sanitizeLabel(label string) string {
	out := make([]byte, len(label))
	for i := range label {
		c := label[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			out[i] = c
		} else {
			out[i] = '_'
		}
	}
	return string(out)
}
