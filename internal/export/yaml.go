package export

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/nicholasgasior/confsnap/internal/baseline"
)

// DriftToYAML writes drift results as a YAML document to w.
func DriftToYAML(w io.Writer, results []baseline.DriftResult, baselineLabel string, capturedAt time.Time) error {
	sorted := make([]baseline.DriftResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	fmt.Fprintf(w, "---\n")
	fmt.Fprintf(w, "baseline: %q\n", baselineLabel)
	fmt.Fprintf(w, "captured_at: %q\n", capturedAt.UTC().Format(time.RFC3339))
	fmt.Fprintf(w, "summary:\n")

	var unchanged, modified, removed int
	for _, r := range sorted {
		switch r.Status {
		case baseline.StatusUnchanged:
			unchanged++
		case baseline.StatusModified:
			modified++
		case baseline.StatusRemoved:
			removed++
		}
	}
	fmt.Fprintf(w, "  unchanged: %d\n", unchanged)
	fmt.Fprintf(w, "  modified: %d\n", modified)
	fmt.Fprintf(w, "  removed: %d\n", removed)
	fmt.Fprintf(w, "files:\n")

	for _, r := range sorted {
		hash := r.CurrentHash
		if hash == "" {
			hash = "-"
		}
		fmt.Fprintf(w, "  - path: %q\n", r.Path)
		fmt.Fprintf(w, "    status: %q\n", strings.ToLower(r.Status.String()))
		fmt.Fprintf(w, "    baseline_hash: %q\n", r.BaselineHash)
		fmt.Fprintf(w, "    current_hash: %q\n", hash)
	}

	return nil
}
