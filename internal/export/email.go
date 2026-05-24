package export

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/confsnap/internal/baseline"
)

// DriftToEmail writes a plain-text email body summarising drift results to w.
// The output is suitable for sending via SMTP or piping to a mail agent.
func DriftToEmail(w io.Writer, results []baseline.DriftResult, snapshotLabel string) error {
	if len(results) == 0 {
		return fmt.Errorf("no drift results provided")
	}

	sorted := make([]baseline.DriftResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "Subject: confsnap drift report — %s\n", snapshotLabel)
	fmt.Fprintf(&buf, "Date: %s\n", time.Now().UTC().Format(time.RFC1123Z))
	fmt.Fprintln(&buf, "MIME-Version: 1.0")
	fmt.Fprintln(&buf, "Content-Type: text/plain; charset=utf-8")
	fmt.Fprintln(&buf)
	fmt.Fprintf(&buf, "Drift Report\n")
	fmt.Fprintf(&buf, "Snapshot : %s\n", snapshotLabel)
	fmt.Fprintf(&buf, "Generated: %s\n", time.Now().UTC().Format(time.RFC3339))
	fmt.Fprintln(&buf)

	// Summary counts
	var added, removed, modified, unchanged int
	for _, r := range sorted {
		switch r.Status {
		case baseline.StatusAdded:
			added++
		case baseline.StatusRemoved:
			removed++
		case baseline.StatusModified:
			modified++
		case baseline.StatusUnchanged:
			unchanged++
		}
	}

	fmt.Fprintf(&buf, "Summary\n")
	fmt.Fprintf(&buf, "  Added    : %d\n", added)
	fmt.Fprintf(&buf, "  Removed  : %d\n", removed)
	fmt.Fprintf(&buf, "  Modified : %d\n", modified)
	fmt.Fprintf(&buf, "  Unchanged: %d\n", unchanged)
	fmt.Fprintln(&buf)

	fmt.Fprintf(&buf, "%-10s  %-60s  %s\n", "STATUS", "PATH", "HASH")
	fmt.Fprintln(&buf, "----------  ------------------------------------------------------------  --------------------------------")

	for _, r := range sorted {
		hash := r.CurrentHash
		if hash == "" {
			hash = "-"
		}
		fmt.Fprintf(&buf, "%-10s  %-60s  %s\n", r.Status.String(), r.Path, hash)
	}

	_, err := io.Copy(w, &buf)
	return err
}
