package diff

import (
	"strings"
	"testing"

	"github.com/user/confsnap/internal/snapshot"
)

func entry(path, hash string, size int64) snapshot.Entry {
	return snapshot.Entry{Path: path, Hash: hash, Size: size}
}

func TestCompare_Added(t *testing.T) {
	before := []snapshot.Entry{entry("/etc/hosts", "aaa", 100)}
	after := []snapshot.Entry{
		entry("/etc/hosts", "aaa", 100),
		entry("/etc/resolv.conf", "bbb", 50),
	}

	r := Compare(before, after)
	if len(r.Added) != 1 || r.Added[0].Path != "/etc/resolv.conf" {
		t.Errorf("expected 1 added file, got %+v", r.Added)
	}
	if len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Errorf("unexpected removals or changes")
	}
}

func TestCompare_Removed(t *testing.T) {
	before := []snapshot.Entry{
		entry("/etc/hosts", "aaa", 100),
		entry("/etc/resolv.conf", "bbb", 50),
	}
	after := []snapshot.Entry{entry("/etc/hosts", "aaa", 100)}

	r := Compare(before, after)
	if len(r.Removed) != 1 || r.Removed[0].Path != "/etc/resolv.conf" {
		t.Errorf("expected 1 removed file, got %+v", r.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	before := []snapshot.Entry{entry("/etc/hosts", "aaa", 100)}
	after := []snapshot.Entry{entry("/etc/hosts", "bbb", 120)}

	r := Compare(before, after)
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed file, got %d", len(r.Changed))
	}
	c := r.Changed[0]
	if c.OldHash != "aaa" || c.NewHash != "bbb" {
		t.Errorf("unexpected hashes: %+v", c)
	}
	if c.OldSize != 100 || c.NewSize != 120 {
		t.Errorf("unexpected sizes: %+v", c)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	before := []snapshot.Entry{entry("/etc/hosts", "aaa", 100)}
	after := []snapshot.Entry{entry("/etc/hosts", "aaa", 100)}

	r := Compare(before, after)
	if len(r.Unchanged) != 1 {
		t.Errorf("expected 1 unchanged file")
	}
	if len(r.Added)+len(r.Removed)+len(r.Changed) != 0 {
		t.Errorf("expected no diffs")
	}
}

func TestFormat_NoChanges(t *testing.T) {
	r := Result{}
	out := Format(r)
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestFormat_WithChanges(t *testing.T) {
	r := Result{
		Added:   []snapshot.Entry{entry("/new", "ccc", 10)},
		Removed: []snapshot.Entry{entry("/old", "ddd", 20)},
		Changed: []Change{{Path: "/etc/hosts", OldHash: "aaa", NewHash: "bbb", OldSize: 100, NewSize: 110}},
	}
	out := Format(r)
	if !strings.Contains(out, "+ /new") {
		t.Errorf("missing added line: %s", out)
	}
	if !strings.Contains(out, "- /old") {
		t.Errorf("missing removed line: %s", out)
	}
	if !strings.Contains(out, "~ /etc/hosts") {
		t.Errorf("missing changed line: %s", out)
	}
}
