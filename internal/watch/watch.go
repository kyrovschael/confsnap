package watch

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// FileEvent represents a change detected in a watched file.
type FileEvent struct {
	Path      string
	ChangedAt time.Time
	OldHash   string
	NewHash   string
}

// Watcher monitors a set of files for changes by comparing SHA-256 hashes.
type Watcher struct {
	paths    []string
	hashes   map[string]string
	Interval time.Duration
}

// New creates a Watcher for the given file paths.
func New(paths []string, interval time.Duration) (*Watcher, error) {
	w := &Watcher{
		paths:    paths,
		hashes:   make(map[string]string),
		Interval: interval,
	}
	for _, p := range paths {
		h, err := hashFile(p)
		if err != nil {
			return nil, fmt.Errorf("watch: initial hash of %s: %w", p, err)
		}
		w.hashes[p] = h
	}
	return w, nil
}

// Poll checks all watched files once and returns any detected FileEvents.
func (w *Watcher) Poll() ([]FileEvent, error) {
	var events []FileEvent
	for _, p := range w.paths {
		newHash, err := hashFile(p)
		if err != nil {
			return nil, fmt.Errorf("watch: polling %s: %w", p, err)
		}
		oldHash := w.hashes[p]
		if newHash != oldHash {
			events = append(events, FileEvent{
				Path:      p,
				ChangedAt: time.Now(),
				OldHash:   oldHash,
				NewHash:   newHash,
			})
			w.hashes[p] = newHash
		}
	}
	return events, nil
}

// hashFile returns the hex-encoded SHA-256 hash of a file's contents.
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
