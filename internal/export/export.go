package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/confsnap/internal/diff"
)

// Record represents a single row in the exported diff output.
type Record struct {
	File      string
	Status    string
	OldHash   string
	NewHash   string
	OldSize   int64
	NewSize   int64
	Timestamp time.Time
}

// ToCSV writes diff results as CSV rows to the provided writer.
// The header row is always written first.
func ToCSV(w io.Writer, results []diff.Result, ts time.Time) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{"file", "status", "old_hash", "new_hash", "old_size", "new_size", "timestamp"}); err != nil {
		return fmt.Errorf("export: write header: %w", err)
	}

	sorted := make([]diff.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].File < sorted[j].File
	})

	for _, r := range sorted {
		row := []string{
			r.File,
			r.Status.String(),
			r.OldHash,
			r.NewHash,
			fmt.Sprintf("%d", r.OldSize),
			fmt.Sprintf("%d", r.NewSize),
			ts.UTC().Format(time.RFC3339),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("export: write row for %q: %w", r.File, err)
		}
	}

	cw.Flush()
	return cw.Error()
}
