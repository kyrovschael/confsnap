package export

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/yourusername/confsnap/internal/baseline"
)

// DriftToMarkdown writes a drift comparison result as a Markdown table to w.
func DriftToMarkdown(w io.Writer, results []baseline.DriftResult, baseline string, at time.Time) error {
	sorted := make([]baseline.DriftResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	_, err := fmt.Fprintf(w, "# Drift Report\n\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "**Baseline:** `%s`  \n**Captured:** `%s`\n\n",
		baseline, at.UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "| File | Status | Expected Hash | Actual Hash |\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "|------|--------|---------------|-------------|\n")
	if err != nil {
		return err
	}

	for _, r := range sorted {
		expected := r.ExpectedHash
		if expected == "" {
			expected = "—"
		}
		actual := r.ActualHash
		if actual == "" {
			actual = "—"
		}
		_, err = fmt.Fprintf(w, "| `%s` | %s | `%s` | `%s` |\n",
			r.Path, r.Status, expected, actual)
		if err != nil {
			return err
		}
	}

	return nil
}
