package export

import (
	"fmt"
	"html"
	"io"
	"sort"
	"time"

	"github.com/yourusername/confsnap/internal/baseline"
)

// DriftToHTML writes a drift comparison report as an HTML table to w.
func DriftToHTML(w io.Writer, results []baseline.DriftResult, label string, capturedAt time.Time) error {
	sorted := make([]baseline.DriftResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Drift Report: %s</title>
<style>
body { font-family: monospace; margin: 2em; }
table { border-collapse: collapse; width: 100%%; }
th, td { border: 1px solid #ccc; padding: 6px 12px; text-align: left; }
th { background: #f0f0f0; }
.added { color: #2a7d2a; }
.removed { color: #b00020; }
.modified { color: #c07000; }
.unchanged { color: #555; }
</style>
</head>
<body>
<h1>Drift Report</h1>
<p><strong>Baseline:</strong> %s</p>
<p><strong>Captured:</strong> %s</p>
<table>
<tr><th>Path</th><th>Status</th><th>Baseline Hash</th><th>Current Hash</th></tr>
`,
		html.EscapeString(label),
		html.EscapeString(label),
		capturedAt.Format(time.RFC3339),
	)

	for _, r := range sorted {
		baseHash := r.BaselineHash
		if baseHash == "" {
			baseHash = "-"
		}
		curHash := r.CurrentHash
		if curHash == "" {
			curHash = "-"
		}
		cssClass := html.EscapeString(r.Status.String())
		fmt.Fprintf(w, "<tr><td>%s</td><td class=\"%s\">%s</td><td>%s</td><td>%s</td></tr>\n",
			html.EscapeString(r.Path),
			cssClass,
			html.EscapeString(r.Status.String()),
			html.EscapeString(baseHash),
			html.EscapeString(curHash),
		)
	}

	fmt.Fprintf(w, "</table>\n</body>\n</html>\n")
	return nil
}
