package schedule

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a scheduled snapshot job.
type Entry struct {
	Label    string        `json:"label"`
	Files    []string      `json:"files"`
	Interval time.Duration `json:"interval"`
	LastRun  time.Time     `json:"last_run,omitempty"`
}

// Schedule holds a collection of scheduled entries.
type Schedule struct {
	Entries []Entry `json:"entries"`
}

func scheduleDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".confsnap", "schedules")
}

// Save persists the schedule to disk.
func Save(s *Schedule) error {
	dir := scheduleDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("schedule: create dir: %w", err)
	}
	path := filepath.Join(dir, "schedule.json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("schedule: create file: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
}

// Load reads the schedule from disk. Returns an empty schedule if not found.
func Load() (*Schedule, error) {
	path := filepath.Join(scheduleDir(), "schedule.json")
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return &Schedule{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("schedule: open: %w", err)
	}
	defer f.Close()
	var s Schedule
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("schedule: decode: %w", err)
	}
	return &s, nil
}

// Add appends or updates an entry by label.
func (s *Schedule) Add(e Entry) {
	for i, existing := range s.Entries {
		if existing.Label == e.Label {
			s.Entries[i] = e
			return
		}
	}
	s.Entries = append(s.Entries, e)
}

// Remove deletes an entry by label. Returns false if not found.
func (s *Schedule) Remove(label string) bool {
	for i, e := range s.Entries {
		if e.Label == label {
			s.Entries = append(s.Entries[:i], s.Entries[i+1:]...)
			return true
		}
	}
	return false
}

// Due returns entries whose next run time has passed.
func (s *Schedule) Due(now time.Time) []Entry {
	var due []Entry
	for _, e := range s.Entries {
		if e.LastRun.IsZero() || now.Sub(e.LastRun) >= e.Interval {
			due = append(due, e)
		}
	}
	return due
}
