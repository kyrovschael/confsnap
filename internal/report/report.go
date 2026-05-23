package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/confsnap/internal/diff"
	"github.com/user/confsnap/internal/snapshot"
)

// Format controls the output format of a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the result of comparing two snapshots.
type Report struct {
	LabelA    string
	LabelB    string
	CreatedAt time.Time
	Entries   []diff.Entry
}

// New builds a Report by comparing two labeled snapshots.
func New(labelA, labelB string) (*Report, error) {
	snapsA, err := snapshot.Load(labelA)
	if err != nil {
		return nil, fmt.Errorf("load snapshot %q: %w", labelA, err)
	}
	snapsB, err := snapshot.Load(labelB)
	if err != nil {
		return nil, fmt.Errorf("load snapshot %q: %w", labelB, err)
	}

	entries := diff.Compare(snapsA, snapsB)
	return &Report{
		LabelA:    labelA,
		LabelB:    labelB,
		CreatedAt: time.Now(),
		Entries:   entries,
	}, nil
}

// Write renders the report to w in the given format.
func (r *Report) Write(w io.Writer, fmt Format) error {
	switch fmt {
	case FormatJSON:
		return writeJSON(w, r)
	default:
		return writeText(w, r)
	}
}

// Summary returns a one-line summary string.
func (r *Report) Summary() string {
	added, removed, changed, unchanged := r.counts()
	parts := []string{
		fmt.Sprintf("+%d added", added),
		fmt.Sprintf("-%d removed", removed),
		fmt.Sprintf("~%d changed", changed),
		fmt.Sprintf("=%d unchanged", unchanged),
	}
	return strings.Join(parts, "  ")
}

func (r *Report) counts() (added, removed, changed, unchanged int) {
	for _, e := range r.Entries {
		switch e.Status {
		case diff.StatusAdded:
			added++
		case diff.StatusRemoved:
			removed++
		case diff.StatusChanged:
			changed++
		case diff.StatusUnchanged:
			unchanged++
		}
	}
	return
}
