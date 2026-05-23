package diff

import (
	"fmt"
	"strings"

	"github.com/user/confsnap/internal/snapshot"
)

// Result holds the comparison between two snapshots.
type Result struct {
	Label  string
	Added  []snapshot.Entry
	Removed []snapshot.Entry
	Changed []Change
	Unchanged []snapshot.Entry
}

// Change represents a file whose content or metadata changed between snapshots.
type Change struct {
	Path    string
	OldHash string
	NewHash string
	OldSize int64
	NewSize int64
}

// Compare diffs two snapshots and returns a Result.
func Compare(before, after []snapshot.Entry) Result {
	beforeMap := make(map[string]snapshot.Entry, len(before))
	for _, e := range before {
		beforeMap[e.Path] = e
	}

	afterMap := make(map[string]snapshot.Entry, len(after))
	for _, e := range after {
		afterMap[e.Path] = e
	}

	var result Result

	for _, e := range after {
		old, exists := beforeMap[e.Path]
		if !exists {
			result.Added = append(result.Added, e)
		} else if old.Hash != e.Hash || old.Size != e.Size {
			result.Changed = append(result.Changed, Change{
				Path:    e.Path,
				OldHash: old.Hash,
				NewHash: e.Hash,
				OldSize: old.Size,
				NewSize: e.Size,
			})
		} else {
			result.Unchanged = append(result.Unchanged, e)
		}
	}

	for _, e := range before {
		if _, exists := afterMap[e.Path]; !exists {
			result.Removed = append(result.Removed, e)
		}
	}

	return result
}

// Format returns a human-readable summary of a diff Result.
func Format(r Result) string {
	var sb strings.Builder

	if len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0 {
		sb.WriteString("No changes detected.\n")
		return sb.String()
	}

	for _, e := range r.Added {
		sb.WriteString(fmt.Sprintf("+ %s (new, %d bytes)\n", e.Path, e.Size))
	}
	for _, e := range r.Removed {
		sb.WriteString(fmt.Sprintf("- %s (removed)\n", e.Path))
	}
	for _, c := range r.Changed {
		sb.WriteString(fmt.Sprintf("~ %s (size %d -> %d, hash %s -> %s)\n",
			c.Path, c.OldSize, c.NewSize,
			truncate(c.OldHash, 8), truncate(c.NewHash, 8)))
	}

	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
