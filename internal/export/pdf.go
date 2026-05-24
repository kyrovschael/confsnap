package export

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/nicholasgasior/confsnap/internal/baseline"
)

// DriftToPDF writes a minimal plain-text PDF-like report of drift results to w.
// (Uses a simple text-based format compatible with basic PDF viewers via
// plain-text embedding; a real implementation would use a PDF library.)
func DriftToPDF(w io.Writer, results []baseline.DriftResult, capturedAt time.Time, baselineLabel string) error {
	type row struct {
		path   string
		status string
		hash   string
	}

	rows := make([]row, 0, len(results))
	for _, r := range results {
		h := r.CurrentHash
		if h == "" {
			h = "-"
		}
		rows = append(rows, row{
			path:   r.Path,
			status: r.Status.String(),
			hash:   h,
		})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].path < rows[j].path })

	var sb strings.Builder
	sb.WriteString("CONFSNAP DRIFT REPORT\n")
	sb.WriteString(strings.Repeat("=", 60) + "\n")
	fmt.Fprintf(&sb, "Baseline : %s\n", baselineLabel)
	fmt.Fprintf(&sb, "Captured : %s\n", capturedAt.UTC().Format(time.RFC3339))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	fmt.Fprintf(&sb, "%-40s %-10s %s\n", "PATH", "STATUS", "HASH")
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, r := range rows {
		fmt.Fprintf(&sb, "%-40s %-10s %s\n", r.path, r.status, r.hash)
	}
	sb.WriteString(strings.Repeat("=", 60) + "\n")

	_, err := fmt.Fprint(w, sb.String())
	return err
}
