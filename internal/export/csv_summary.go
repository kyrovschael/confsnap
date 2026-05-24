package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/yourusername/confsnap/internal/baseline"
)

// DriftSummaryRow represents a single row in the drift summary CSV export.
type DriftSummaryRow struct {
	File      string
	Status    string
	Baseline  string
	CheckedAt time.Time
}

// DriftToCSV writes a slice of DriftResult values to w in CSV format.
// Rows are sorted by file path for deterministic output.
func DriftToCSV(w io.Writer, results []baseline.DriftResult, checkedAt time.Time, baselineName string) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{"file", "status", "baseline", "checked_at"}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	sorted := make([]baseline.DriftResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].File < sorted[j].File
	})

	for _, r := range sorted {
		row := []string{
			r.File,
			r.Status.String(),
			baselineName,
			checkedAt.UTC().Format(time.RFC3339),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
