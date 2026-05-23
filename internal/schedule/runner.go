package schedule

import (
	"fmt"
	"io"
	"time"

	"github.com/user/confsnap/internal/snapshot"
	"github.com/user/confsnap/internal/snapshot/store"
)

// RunOnce checks for due entries, captures snapshots, and updates LastRun.
// It writes a summary line per processed entry to w.
func RunOnce(w io.Writer, now time.Time) error {
	s, err := Load()
	if err != nil {
		return fmt.Errorf("runner: load schedule: %w", err)
	}
	due := s.Due(now)
	if len(due) == 0 {
		fmt.Fprintln(w, "no entries due")
		return nil
	}
	for i, e := range due {
		var snaps []snapshot.Snapshot
		for _, path := range e.Files {
			snap, err := snapshot.Capture(path)
			if err != nil {
				fmt.Fprintf(w, "warn: capture %s: %v\n", path, err)
				continue
			}
			snaps = append(snaps, snap)
		}
		if len(snaps) == 0 {
			continue
		}
		if err := store.Save(e.Label, snaps); err != nil {
			return fmt.Errorf("runner: save snapshot for %s: %w", e.Label, err)
		}
		fmt.Fprintf(w, "snapshotted %s (%d files)\n", e.Label, len(snaps))
		// Update LastRun on the matching schedule entry.
		for j := range s.Entries {
			if s.Entries[j].Label == due[i].Label {
				s.Entries[j].LastRun = now
			}
		}
	}
	return Save(s)
}
