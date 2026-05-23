package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Baseline represents a named reference point for configuration state.
type Baseline struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Files     []string          `json:"files"`
	Hashes    map[string]string `json:"hashes"`
}

const baselineDir = ".confsnap/baselines"

// Save writes the baseline to disk as a JSON file.
func Save(b *Baseline) error {
	if err := os.MkdirAll(baselineDir, 0o755); err != nil {
		return fmt.Errorf("create baseline dir: %w", err)
	}
	path := filepath.Join(baselineDir, b.Name+".json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create baseline file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// Load reads a baseline by name from disk.
func Load(name string) (*Baseline, error) {
	path := filepath.Join(baselineDir, name+".json")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline %q not found", name)
		}
		return nil, fmt.Errorf("open baseline: %w", err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("decode baseline: %w", err)
	}
	return &b, nil
}

// List returns the names of all saved baselines.
func List() ([]string, error) {
	entries, err := os.ReadDir(baselineDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read baseline dir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
