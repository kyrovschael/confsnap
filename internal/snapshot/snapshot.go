package snapshot

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// Snapshot represents the state of a single config file at a point in time.
type Snapshot struct {
	Path      string    `json:"path"`
	Hash      string    `json:"hash"`
	Size      int64     `json:"size"`
	CapturedAt time.Time `json:"captured_at"`
	Content   []byte    `json:"content,omitempty"`
}

// Capture reads the file at the given path and returns a Snapshot.
func Capture(filePath string, includeContent bool) (*Snapshot, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", filePath, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file %q: %w", filePath, err)
	}

	h := sha256.New()
	var content []byte

	if includeContent {
		content, err = io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("reading file %q: %w", filePath, err)
		}
		h.Write(content)
	} else {
		if _, err := io.Copy(h, f); err != nil {
			return nil, fmt.Errorf("hashing file %q: %w", filePath, err)
		}
	}

	return &Snapshot{
		Path:       filePath,
		Hash:       fmt.Sprintf("%x", h.Sum(nil)),
		Size:       info.Size(),
		CapturedAt: time.Now().UTC(),
		Content:    content,
	}, nil
}
