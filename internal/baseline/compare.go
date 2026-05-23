package baseline

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Drift describes how a single file has changed relative to the baseline.
type Drift struct {
	Path    string
	Status  DriftStatus
	OldHash string
	NewHash string
}

// DriftStatus categorises a file's change state.
type DriftStatus int

const (
	DriftUnchanged DriftStatus = iota
	DriftModified
	DriftAdded
	DriftRemoved
)

func (s DriftStatus) String() string {
	switch s {
	case DriftUnchanged:
		return "unchanged"
	case DriftModified:
		return "modified"
	case DriftAdded:
		return "added"
	case DriftRemoved:
		return "removed"
	default:
		return "unknown"
	}
}

// CheckDrift compares the baseline hashes against the current state of files on disk.
func CheckDrift(b *Baseline) ([]Drift, error) {
	seen := make(map[string]bool)
	var drifts []Drift

	for _, path := range b.Files {
		oldHash := b.Hashes[path]
		newHash, err := hashFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				drifts = append(drifts, Drift{Path: path, Status: DriftRemoved, OldHash: oldHash})
				seen[path] = true
				continue
			}
			return nil, fmt.Errorf("hash %q: %w", path, err)
		}
		seen[path] = true
		if newHash == oldHash {
			drifts = append(drifts, Drift{Path: path, Status: DriftUnchanged, OldHash: oldHash, NewHash: newHash})
		} else {
			drifts = append(drifts, Drift{Path: path, Status: DriftModified, OldHash: oldHash, NewHash: newHash})
		}
	}

	return drifts, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
